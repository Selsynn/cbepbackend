package interactiondiscordimpl

import (
	"fmt"

	"github.com/Selsynn/cbepbackend/business/command"
	"github.com/Selsynn/cbepbackend/business/player"
	"github.com/Selsynn/cbepbackend/business/town"
	"github.com/Selsynn/cbepbackend/business/user"
	"github.com/Selsynn/cbepbackend/communication"
	"github.com/Selsynn/cbepbackend/discord"
	"github.com/Selsynn/cbepbackend/discord/discordreaction"
	"github.com/Selsynn/cbepbackend/talker"
	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
)

type InteractionDiscord struct {
	Servers map[discord.ServerID]*discord.Server
}

func (i *InteractionDiscord) GetOrCreateServer(serverID discord.ServerID, channelID discord.ChannelID) *discord.Server {
	if _, ok := i.Servers[serverID]; !ok {
		newTown := town.New()
		i.Servers[serverID] = &discord.Server{
			ID:                       serverID,
			ChannelID:                channelID,
			PlayerAdventurers:        make(map[user.ID]player.ID),
			TownID:                   newTown.ID,
			WaitingActionsForPlayers: []communication.ActionFromManager{},
		}
	}
	return i.Servers[serverID]
}

func (i *InteractionDiscord) GetTown(server *discord.Server) town.ID {
	return server.TownID
}

func (i *InteractionDiscord) GetActionToManager(message talker.MessageReceived) communication.ActionToManager {
	result := communication.ActionToManager{}

	switch v := message.(type) {
	case *discord.TextReceiveDiscord:
		server := i.GetOrCreateServer(v.Server.ID, v.Server.ChannelID)
		result.TownID = i.GetTown(server)
		result.PlayerID = i.GetPlayer(v.User, server)
		result.ActionID = v.Message
		result.Command = i.GetCommandFromText(v.Text)
	case *discord.ReactionReceiveDiscord:
		server := i.GetOrCreateServer(v.Server.ID, v.Server.ChannelID)
		result.TownID = i.GetTown(server)
		result.PlayerID = i.GetPlayer(v.User, server)
		result.ActionID = v.Message
		result.Command = i.GetCommandFromReaction(v.Reaction)
	}

	return result
}

func (i *InteractionDiscord) GetActionFromManager(message communication.ActionFromManager) talker.MessageSent {
	server := i.GetServerFromTown(message.TownID)
	result := discord.MessageSentDiscord{
		Server: discord.ServerDiscord{
			ID:        server.ID,
			ChannelID: server.ChannelID,
		},
	}
	for command := range message.Callback {
		result.ReactionIDs = append(result.ReactionIDs, discordreaction.Match2Reaction(command))
	}
	result.Text = discordgo.MessageEmbed{
		Title:       "You have a new message",
		Description: message.Content.Text,
	}

	return &result
}

func (i *InteractionDiscord) GetPlayer(user user.ID, server *discord.Server) player.ID {
	if _, ok := server.PlayerAdventurers[user]; !ok {
		server.PlayerAdventurers[user] = player.ID(uuid.New().String())
	}
	return server.PlayerAdventurers[user]
}

func (i *InteractionDiscord) GetCommandFromReaction(reaction discordreaction.ID) command.Command {
	return command.CommandSimple{
		Id: discordreaction.Match2Command(reaction),
	}
}

func (i *InteractionDiscord) GetCommandFromText(text string) command.Command {
	cmd, err := command.Parse(text)
	if err != nil {
		fmt.Println(err.Error())
	}
	return cmd
}

func (i *InteractionDiscord) GetCallback(toManager communication.ActionToManager) func() *communication.ActionFromManager {
	server := i.GetServerFromTown(toManager.TownID)

	checkInAllowList := func(allowList []*player.ID, playerID player.ID) bool {
		for _, allowed := range allowList {
			if *allowed == playerID {
				return true
			}
		}
		return false
	}

	for _, waitingAction := range server.WaitingActionsForPlayers {
		if waitingAction.MessageID == toManager.ActionID && checkInAllowList(waitingAction.AllowList, toManager.PlayerID) {
			for expected, callback := range waitingAction.Callback {
				if expected == toManager.Command.ID() {
					return callback
				}
			}
		}
	}
	return nil
}

func (i *InteractionDiscord) GetServerFromTown(town town.ID) *discord.Server {
	for _, server := range i.Servers {
		if server.TownID == town {
			return server
		}
	}
	panic("No server found for the town " + town)

}

func (i *InteractionDiscord) AddCallback(fromManager communication.ActionFromManager, actionID communication.ActionID) {
	server := i.GetServerFromTown(fromManager.TownID)
	fromManager.MessageID = actionID
	server.WaitingActionsForPlayers = append(server.WaitingActionsForPlayers, fromManager)
}

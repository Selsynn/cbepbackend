package talkerdiscord

import (
	"fmt"

	"github.com/Selsynn/craft-build-explore-protect-backend/business/user"
	"github.com/Selsynn/craft-build-explore-protect-backend/discord"
	"github.com/Selsynn/craft-build-explore-protect-backend/discord/discordreaction"
	"github.com/Selsynn/craft-build-explore-protect-backend/talker"
	"github.com/bwmarrin/discordgo"
)

type TalkerDiscord struct {
	session   *discordgo.Session
	messageCh chan talker.MessageReceived
	//Servers   []*ServerDiscord
}

func NewTalkerDiscord(token string) (impl *TalkerDiscord, shutdown func()) {
	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		panic("error creating Discord session," + err.Error())
	}

	t := &TalkerDiscord{
		session:   dg,
		messageCh: make(chan talker.MessageReceived, 1),
		//Servers:   make([]*ServerDiscord, 0),
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate(t))

	// Register the ReactionGet all
	dg.AddHandler(messageReactionAdd(t))

	// dg.AddHandler(func(s *discordgo.Session, m *discordgo.MessageReactionRemove) {
	// 	fmt.Printf("Remove %#v\n", m.MessageReaction)
	// })
	// dg.AddHandler(func(s *discordgo.Session, m *discordgo.MessageReactionRemoveAll) {
	// 	fmt.Printf("Remove all %#v\n", m.MessageReaction.Emoji)
	// })

	// In this example, we only care about receiving message events.
	dg.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuildMessages | discordgo.IntentsGuildMessageReactions)

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		panic("error opening connection," + err.Error())
	}

	return t, t.Shutdown
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func messageCreate(t *TalkerDiscord) func(s *discordgo.Session, m *discordgo.MessageCreate) {
	return func(s *discordgo.Session, m *discordgo.MessageCreate) {
		// Ignore all messages created by the bot itself
		// This isn't required in this specific example but it's a good practice.
		if m.Author.ID == s.State.User.ID {
			return
		}

		mess := &discord.TextReceiveDiscord{
			Server: discord.ServerDiscord{
				ChannelID: discord.ChannelID(m.ChannelID),
				ID:        discord.ServerID(m.GuildID),
			},
			Message: discord.MessageID(m.Message.ID),
			Text:    m.Content,
			User:    user.ID(m.Author.ID),
		}

		t.messageCh <- mess
	}

}

func messageReactionAdd(t *TalkerDiscord) func(s *discordgo.Session, m *discordgo.MessageReactionAdd) {
	return func(s *discordgo.Session, m *discordgo.MessageReactionAdd) {
		if m.UserID == s.State.User.ID {
			return
		}

		fmt.Printf("Add %#v\n", m.MessageReaction)
		mess := &discord.ReactionReceiveDiscord{
			Server: discord.ServerDiscord{
				ChannelID: discord.ChannelID(m.ChannelID),
				ID:        discord.ServerID(m.GuildID),
			},
			Message:  discord.MessageID(m.MessageID),
			Reaction: discordreaction.ID(m.Emoji.Name),
			User:     user.ID(m.UserID),
		}

		t.messageCh <- mess
	}
}

func (t TalkerDiscord) Read() chan talker.MessageReceived {
	return t.messageCh
}

func (t TalkerDiscord) Write(i talker.MessageSent) {
	m := i.(*discord.MessageSentDiscord)

	message, err := t.session.ChannelMessageSendEmbed(string(m.Server.ChannelID), &m.Text)

	if err != nil {
		panic(err.Error())
	}

	for _, reaction := range m.ReactionIDs {
		err = t.session.MessageReactionAdd(string(m.Server.ChannelID), message.ID, string(reaction))

		if err != nil {
			fmt.Printf("Trying to send %s %s\n", reaction, err)
		}
	}
}

func (t TalkerDiscord) Shutdown() {
	// Cleanly close down the Discord session.
	t.session.Close()
}

// func (t *TalkerDiscord) FindOrCreateServer(channelId string) *ServerDiscord {
// 	for _, item := range t.Servers {
// 		if item.GetId() == channelId {
// 			return item
// 		}
// 	}

// 	//we didn't find anything, time to create it
// 	s := &ServerDiscord{
// 		channelId: channelId,
// 		name:      "Basic Name",
// 	}
// 	t.Servers = append(t.Servers, s)
// 	return s
// }

// func (t TalkerDiscord) GetServers() []Server {
// 	return t.Servers
// }

// func (s ServerDiscord) GetName() string {
// 	return s.name
// }

// func (s ServerDiscord) GetId() string {
// 	return s.channelId
// }

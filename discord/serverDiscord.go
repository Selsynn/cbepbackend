package discord

import (
	"github.com/Selsynn/DiscordBotTest1/business/player"
	"github.com/Selsynn/DiscordBotTest1/business/town"
	"github.com/Selsynn/DiscordBotTest1/business/user"
	"github.com/Selsynn/DiscordBotTest1/communication"
)

type ServerID string
type ChannelID string

type Server struct {
	ID                       ServerID
	ChannelID                ChannelID
	PlayerAdventurers        map[user.ID]player.ID
	WaitingActionsForPlayers []communication.ActionFromManager
	TownID                   town.ID
}



package discord

import (
	"github.com/Selsynn/cbepbackend/business/player"
	"github.com/Selsynn/cbepbackend/business/town"
	"github.com/Selsynn/cbepbackend/business/user"
	"github.com/Selsynn/cbepbackend/communication"
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

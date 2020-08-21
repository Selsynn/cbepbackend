package discord

import (
	"github.com/Selsynn/craft-build-explore-protect-backend/business/player"
	"github.com/Selsynn/craft-build-explore-protect-backend/business/town"
	"github.com/Selsynn/craft-build-explore-protect-backend/business/user"
	"github.com/Selsynn/craft-build-explore-protect-backend/communication"
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

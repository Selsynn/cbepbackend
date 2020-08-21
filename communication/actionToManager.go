package communication

import (
	"github.com/Selsynn/craft-build-explore-protect-backend/business/command"
	"github.com/Selsynn/craft-build-explore-protect-backend/business/player"
	"github.com/Selsynn/craft-build-explore-protect-backend/business/town"
)

type ActionToManager struct {
	TownID   town.ID
	PlayerID player.ID
	Command  command.Command
}

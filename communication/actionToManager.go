package communication

import (
	"github.com/Selsynn/cbepbackend/business/command"
	"github.com/Selsynn/cbepbackend/business/player"
	"github.com/Selsynn/cbepbackend/business/town"
)

type ActionToManager struct {
	TownID   town.ID
	PlayerID player.ID
	ActionID *ActionID
	ParentID *ActionID
	Command  command.Command
}

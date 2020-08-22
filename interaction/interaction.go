package interaction

import (
	"github.com/Selsynn/cbepbackend/business/town"
	"github.com/Selsynn/cbepbackend/communication"
	"github.com/Selsynn/cbepbackend/talker"
)

type Interaction interface {
	GetActionToManager(message talker.MessageReceived, createTown func() town.ID) communication.ActionToManager
	GetActionFromManager(message communication.ActionFromManager) talker.MessageSent
	GetCallback(toManager communication.ActionToManager) func(communication.ActionToManager) *communication.ActionFromManager
	AddCallback(fromManager communication.ActionFromManager, actionID communication.ActionID)
	CleanCallback(toManager communication.ActionToManager)
}

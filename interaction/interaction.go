package interaction

import (
	"github.com/Selsynn/cbepbackend/communication"
	"github.com/Selsynn/cbepbackend/talker"
)

type Interaction interface {
	GetActionToManager(message talker.MessageReceived) communication.ActionToManager
	GetActionFromManager(message communication.ActionFromManager) talker.MessageSent
	GetCallback(toManager communication.ActionToManager) func() *communication.ActionFromManager
	AddCallback(fromManager communication.ActionFromManager, actionID communication.ActionID)
}

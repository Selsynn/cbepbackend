package interaction

import (
	"github.com/Selsynn/craft-build-explore-protect-backend/communication"
	"github.com/Selsynn/craft-build-explore-protect-backend/talker"
)

type Interaction interface {
	GetActionToManager(message talker.MessageReceived) communication.ActionToManager
	GetActionFromManager(message communication.ActionFromManager) talker.MessageSent
}

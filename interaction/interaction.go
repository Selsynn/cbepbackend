package interaction

import (
	"github.com/Selsynn/DiscordBotTest1/communication"
	"github.com/Selsynn/DiscordBotTest1/talker"
)

type Interaction interface {
	GetActionToManager(message talker.MessageReceived) communication.ActionToManager
	GetActionFromManager(message communication.ActionFromManager) talker.MessageSent
}

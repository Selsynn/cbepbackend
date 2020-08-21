package communication

import (
	"github.com/Selsynn/DiscordBotTest1/business/command"
	"github.com/Selsynn/DiscordBotTest1/business/player"
	"github.com/Selsynn/DiscordBotTest1/business/town"
)

type ActionToManager struct {
	TownID   town.ID
	PlayerID player.ID
	Command  command.Command
}

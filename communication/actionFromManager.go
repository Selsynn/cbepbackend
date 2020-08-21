package communication

import (
	"github.com/Selsynn/DiscordBotTest1/business/command"
	"github.com/Selsynn/DiscordBotTest1/business/player"
	"github.com/Selsynn/DiscordBotTest1/business/town"
)

type ActionFromManager struct {
	Parent    *ActionFromManager
	Content   ContentMessage
	AllowList []*player.ID
	TownID    town.ID
	Callback  map[command.Command]func(ActionFromManager)
	//CleanUp  func()
}

type ContentMessage struct {
	Text       string
	//AllActions []*ActionFromManager
}

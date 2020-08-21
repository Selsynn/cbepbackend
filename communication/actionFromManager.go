package communication

import (
	"github.com/Selsynn/cbepbackend/business/command"
	"github.com/Selsynn/cbepbackend/business/player"
	"github.com/Selsynn/cbepbackend/business/town"
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
	Text string
	//AllActions []*ActionFromManager
}

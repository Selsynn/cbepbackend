package communication

import (
	"github.com/Selsynn/cbepbackend/business/command"
	"github.com/Selsynn/cbepbackend/business/player"
	"github.com/Selsynn/cbepbackend/business/town"
)

type ActionID string

type ActionFromManager struct {
	MessageID ActionID
	Parent    *ActionID
	Content   ContentMessage
	AllowList []*player.ID
	TownID    town.ID
	Callback  map[command.ID]func(ActionToManager) *ActionFromManager
	//CleanUp  func()
}

type ContentMessage struct {
	Text       string
	ActionFlag map[command.ID]string
	//AllActions []*ActionFromManager
}

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
	Callback  map[command.ID]ActionCallback
	//CleanUp  func()
}

type ContentMessage struct {
	Text       string
	ActionFlag map[command.ID]string
}

type DescriptionAction struct {
	CID         command.ID
	Description string
	Callback    ActionCallback
}

type ActionCallback func(ActionToManager) *ActionFromManager

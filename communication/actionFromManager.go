package communication

import (
	"github.com/Selsynn/craft-build-explore-protect-backend/business/command"
	"github.com/Selsynn/craft-build-explore-protect-backend/business/player"
	"github.com/Selsynn/craft-build-explore-protect-backend/business/town"
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

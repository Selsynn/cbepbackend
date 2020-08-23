package manager

import (
	"github.com/Selsynn/cbepbackend/business/characters"
	"github.com/Selsynn/cbepbackend/business/command"
	"github.com/Selsynn/cbepbackend/business/player"
	"github.com/Selsynn/cbepbackend/business/town"
	"github.com/Selsynn/cbepbackend/communication"
)

func (m *Manager) Craft(message communication.ActionToManager) *communication.ActionFromManager {
	callbacks := map[command.ID]communication.DescriptionAction{
		command.Profile: m.ProfileCallback(),
	}

	return &communication.ActionFromManager{
		TownID: message.TownID,
		AllowList: []*player.ID{
			&message.PlayerID,
		},
		Callback: callbacks,
		Content: communication.ContentMessage{
			Text: "Craft (WIP)",
		},
		Parent: message.ParentID,
	}
}

type CraftActionRegistry struct {
	ID             characters.NPCName
	CID            command.ID
	Description    string
	ResultCallback string
	sideEffects    func(t *town.Town)
}

var craftKnowledgeList map[characters.NPCName]map[characters.RelationLevel][]*CraftActionRegistry

func init() {
	// weapon1 := []*CraftActionRegistry{
	// 	// {
	// 	// 	ID: ,
	// 	// },
	// }
}

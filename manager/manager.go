package manager

import (
	"fmt"

	"github.com/Selsynn/cbepbackend/business/command"
	"github.com/Selsynn/cbepbackend/business/player"
	"github.com/Selsynn/cbepbackend/business/town"
	"github.com/Selsynn/cbepbackend/communication"
)

type Manager struct {
	Towns map[town.ID]*town.Town
}

func (m *Manager) Process(message communication.ActionToManager) *communication.ActionFromManager {
	//Get town
	// t, found := m.Towns[message.TownID]

	// if !found {
	// 	//No city found
	// 	t = town.New()
	// 	m.Towns[message.TownID] = t
	// 	fmt.Printf("Creation of a new city \n%#v\n", t)
	// }

	//a := t.Adventurers[0]

	fmt.Printf("Describe WORLD \n%#v\n", m.Towns)

	result := &communication.ActionFromManager{
		TownID: message.TownID,
		AllowList: []*player.ID{
			&message.PlayerID,
		},
	}

	var resultContent string
	resultActionFlag := make(map[command.ID]string)

	switch message.Command.ID() {
	case command.Explore:
		return m.Explore(message)
	case command.Profile:
		return m.Profile(message)

	// case command.ViewShop:
	// 	resultContent = t.DescribeCity()
	// case command.NewMerchant:
	// 	t.CreateMerchant()
	// 	resultContent = "Merchant created." + t.DescribeCity()
	// 	result.Callback = map[command.ID]func() *communication.ActionFromManager{
	// 		command.Accept: func() *communication.ActionFromManager {
	// 			fmt.Println("Accepted !")
	// 			return nil
	// 		},
	// 	}
	// case command.ViewHeros:
	// 	resultContent = "You asked to see the heros" + t.DescribeHeros()
	// case command.Craft:
	// 	craft := message.Command.(command.CommandCraft)
	// 	t.Craft(t.Crafters[0], t.GetItem(craft.ItemID), a)
	default:
		resultContent = fmt.Sprintf("Il n'y a rien a cette adresse. List of all the command currently supported: **%v**", command.ListAll())
	}

	result.Content = communication.ContentMessage{
		Text:       resultContent,
		ActionFlag: resultActionFlag,
	}

	return result
}

func NewManager() Manager {
	return Manager{
		Towns: make(map[town.ID]*town.Town),
	}
}

func (m *Manager) Explore(message communication.ActionToManager) *communication.ActionFromManager {
	return &communication.ActionFromManager{
		TownID: message.TownID,
		AllowList: []*player.ID{
			&message.PlayerID,
		},
		Callback: map[command.ID]func(communication.ActionToManager) *communication.ActionFromManager{
			command.Wood: func(action communication.ActionToManager) *communication.ActionFromManager {
				return m.ExploreWoods(action)
			},
		},
		Content: communication.ContentMessage{
			Text: "Exploration asked.",
			ActionFlag: map[command.ID]string{
				command.Wood: "Go into the Enchanted Forest",
			},
		},
		Parent: message.ParentID,
	}
}

func (m *Manager) Profile(message communication.ActionToManager) *communication.ActionFromManager {
	return &communication.ActionFromManager{
		TownID: message.TownID,
		AllowList: []*player.ID{
			&message.PlayerID,
		},
		Callback: map[command.ID]func(communication.ActionToManager) *communication.ActionFromManager{
			command.Explore: func(action communication.ActionToManager) *communication.ActionFromManager {
				return m.Explore(action)
			},
			command.Craft: func(action communication.ActionToManager) *communication.ActionFromManager {
				return m.Craft(action)
			},
			command.Sell: func(action communication.ActionToManager) *communication.ActionFromManager {
				return m.Sell(action)
			},
		},
		Content: communication.ContentMessage{
			Text: "Profile asked.",
			ActionFlag: map[command.ID]string{
				command.Explore: "Go into exporation",
				command.Craft:   "Craft something",
				command.Sell:    "Sell something",
			},
		},
		Parent: message.ParentID,
	}
}

func (m *Manager) Craft(message communication.ActionToManager) *communication.ActionFromManager {
	return &communication.ActionFromManager{
		TownID: message.TownID,
		AllowList: []*player.ID{
			&message.PlayerID,
		},
		Callback: map[command.ID]func(communication.ActionToManager) *communication.ActionFromManager{
			command.Explore: func(action communication.ActionToManager) *communication.ActionFromManager {
				return m.Explore(action)
			},
			command.Craft: func(action communication.ActionToManager) *communication.ActionFromManager {
				return m.Craft(action)
			},
			command.Sell: func(action communication.ActionToManager) *communication.ActionFromManager {
				return m.Sell(action)
			},
		},
		Content: communication.ContentMessage{
			Text: "Craft (WIP)",
			ActionFlag: map[command.ID]string{
				command.Explore: "Go into exporation",
				command.Craft:   "Craft something",
				command.Sell:    "Sell something",
			},
		},
		Parent: message.ParentID,
	}
}

func (m *Manager) Sell(message communication.ActionToManager) *communication.ActionFromManager {
	return &communication.ActionFromManager{
		TownID: message.TownID,
		AllowList: []*player.ID{
			&message.PlayerID,
		},
		Callback: map[command.ID]func(communication.ActionToManager) *communication.ActionFromManager{
			command.Explore: func(action communication.ActionToManager) *communication.ActionFromManager {
				return m.Explore(action)
			},
			command.Craft: func(action communication.ActionToManager) *communication.ActionFromManager {
				return m.Craft(action)
			},
			command.Sell: func(action communication.ActionToManager) *communication.ActionFromManager {
				return m.Sell(action)
			},
		},
		Content: communication.ContentMessage{
			Text: "Sell (WIP)",
			ActionFlag: map[command.ID]string{
				command.Explore: "Go into exporation",
				command.Craft:   "Craft something",
				command.Sell:    "Sell something",
			},
		},
		Parent: message.ParentID,
	}
}

func (m *Manager) ExploreWoods(message communication.ActionToManager) *communication.ActionFromManager {
	fmt.Println("Explore - Enchanted Forest !")
	resultExploration := "You have found xxx"
	return &communication.ActionFromManager{
		TownID: message.TownID,
		AllowList: []*player.ID{
			&message.PlayerID,
		},
		Parent: message.ParentID,
		Content: communication.ContentMessage{
			Text: resultExploration,
			ActionFlag: map[command.ID]string{
				command.Accept: "Accept",
				command.Refuse: "Refuse (try a reload)",
			},
		},
		Callback: map[command.ID]func(communication.ActionToManager) *communication.ActionFromManager{
			command.Accept: func(action communication.ActionToManager) *communication.ActionFromManager {
				fmt.Println("Explore - Enchanted Forest - Accepted !")
				return m.Profile(action)
			},
			command.Refuse: func(action communication.ActionToManager) *communication.ActionFromManager {
				fmt.Println("Explore - Enchanted Forest - Retry !")
				return m.ExploreWoods(action)
			},
		},
	}
}

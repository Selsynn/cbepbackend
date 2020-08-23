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
	t, found := m.Towns[message.TownID]

	if !found {
		//No city found
		panic("No town found!")
	}

	fmt.Printf("Describe TOWN \n%#v\n", t)

	_, AdventurerFound := t.Adventurers[message.PlayerID]

	if !AdventurerFound && message.Command.ID() != command.CreateProfile {
		//This adventurer is not known.
		return &communication.ActionFromManager{
			Content: communication.ContentMessage{
				Text: "Please create a Adventurer by running the `create profile: YOURNAME` (no funny character in it => no space, no ponctuation)",
			},
			TownID: t.ID,
		}
	}

	result := &communication.ActionFromManager{
		TownID: message.TownID,
		AllowList: []*player.ID{
			&message.PlayerID,
		},
	}

	switch message.Command.ID() {
	case command.Explore:
		return m.Explore(message)
	case command.Profile:
		return m.Profile(message)
	case command.CreateProfile:
		if AdventurerFound {
			return &communication.ActionFromManager{
				Content: communication.ContentMessage{
					Text: "You already have a Adventurer! You can't create another!",
				},
				TownID: t.ID,
				Callback: map[command.ID]communication.DescriptionAction{
					command.Profile: m.ProfileCallback(),
				},
			}
		}
		createProfile := message.Command.(command.CommandCreateName)
		t.CreateAdventurer(createProfile.Name, message.PlayerID)
		result.Content = communication.ContentMessage{
			Text: "Your character is created! Welcome " + createProfile.Name,
		}
		result.Callback = map[command.ID]communication.DescriptionAction{
			command.Profile: m.ProfileCallback(),
		}

	default:
		result.Content = communication.ContentMessage{
			Text: fmt.Sprintf("Il n'y a rien a cette adresse. List of all the command currently supported: **%v**", command.ListAll()),
		}
	}

	return result
}

func NewManager() Manager {
	return Manager{
		Towns: make(map[town.ID]*town.Town),
	}
}

func (m *Manager) CreateTown() town.ID {
	t := town.New()
	m.Towns[t.ID] = t
	fmt.Printf("Creation of a new city \n%#v\n", t)
	return t.ID
}

func (m *Manager) Profile(message communication.ActionToManager) *communication.ActionFromManager {
	a := m.Towns[message.TownID].Adventurers[message.PlayerID]

	itemsDescr := ""
	for _, item := range a.Items {
		itemsDescr += item.Describe() + " "
	}
	return &communication.ActionFromManager{
		TownID: message.TownID,
		AllowList: []*player.ID{
			&message.PlayerID,
		},
		Callback: map[command.ID]communication.DescriptionAction{
			command.Explore: m.ExploreCallback(),
			command.Craft:   m.CraftCallback(),
			command.Sell:    m.SellCallback(),
		},
		Content: communication.ContentMessage{
			Text: "Profile asked.",
			OtherField: map[string]string{
				"Name":      a.Name,
				"Inventory": itemsDescr,
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
		Callback: map[command.ID]communication.DescriptionAction{
			command.Explore: m.ExploreCallback(),
			command.Craft:   m.CraftCallback(),
			command.Sell:    m.SellCallback(),
		},
		Content: communication.ContentMessage{
			Text: "Sell (WIP)",
		},
		Parent: message.ParentID,
	}
}

// func (m *Manager) CraftElement(message communication.ActionToManager, specialty town.Profession) *communication.ActionFromManager {
// 	fmt.Println("Craft - " + specialty)
// 	resultCraft := "You asked the " + string(specialty) + " to craft something."
// 	result := &communication.ActionFromManager{
// 		TownID: message.TownID,
// 		AllowList: []*player.ID{
// 			&message.PlayerID,
// 		},
// 		Parent: message.ParentID,
// 		Content: communication.ContentMessage{
// 			Text: resultCraft,
// 			ActionFlag: map[command.ID]string{
// 				//command.Bow:    "Craft a bow (cost 10 woods - time 1 hu)",
// 				command.Refuse: "Cancel",
// 			},
// 		},
// 		Callback: map[command.ID]func(communication.ActionToManager) *communication.ActionFromManager{
// 			// command.Bow: func(action communication.ActionToManager) *communication.ActionFromManager {
// 			// 	fmt.Println("Craft bow")
// 			// 	return m.Profile(action)
// 			// },
// 			command.Refuse: func(action communication.ActionToManager) *communication.ActionFromManager {
// 				fmt.Println("Cancel ")
// 				return m.Craft(action)
// 			},
// 		},
// 	}

// 	crafter := m.Towns[message.TownID].NPC[specialty]

// 	for _, element := range crafter.Knowledge {
// 		stringDescr := ""
// 		for _, cost := range element.Cost {
// 			stringDescr += cost.Describe() + " "
// 		}
// 		result.Content.Text = fmt.Sprintf("%s \n***%s***\n--Cost: %s\n--Time: %f Time Unit", result.Content.Text, element.Item.Name, stringDescr, element.WorkCost)
// 	}

// 	return result
// }

// func (m *Manager) ExploreWoods(message communication.ActionToManager) *communication.ActionFromManager {
// 	fmt.Println("Explore - Enchanted Forest !")
// 	resultExploration := "You have found xxx"
// 	return &communication.ActionFromManager{
// 		TownID: message.TownID,
// 		AllowList: []*player.ID{
// 			&message.PlayerID,
// 		},
// 		Parent: message.ParentID,
// 		Content: communication.ContentMessage{
// 			Text: resultExploration,
// 			ActionFlag: map[command.ID]string{
// 				command.Accept: "Accept",
// 				command.Refuse: "Refuse (try a reload)",
// 			},
// 		},
// 		Callback: map[command.ID]func(communication.ActionToManager) *communication.ActionFromManager{
// 			command.Accept: func(action communication.ActionToManager) *communication.ActionFromManager {
// 				fmt.Println("Explore - Enchanted Forest - Accepted !")
// 				return m.Profile(action)
// 			},
// 			command.Refuse: func(action communication.ActionToManager) *communication.ActionFromManager {
// 				fmt.Println("Explore - Enchanted Forest - Retry !")
// 				return m.ExploreWoods(action)
// 			},
// 		},
// 	}
// }

func (m *Manager) ProfileCallback() communication.DescriptionAction {
	return communication.DescriptionAction{
		CID:         command.Profile,
		Callback:    m.Profile,
		Description: "See my profile",
	}
}

func (m *Manager) ExploreCallback() communication.DescriptionAction {
	return communication.DescriptionAction{
		CID:         command.Explore,
		Callback:    m.Explore,
		Description: "Explore",
	}
}

func (m *Manager) CraftCallback() communication.DescriptionAction {
	return communication.DescriptionAction{
		CID:         command.Craft,
		Callback:    m.Craft,
		Description: "Craft",
	}
}

func (m *Manager) SellCallback() communication.DescriptionAction {
	return communication.DescriptionAction{
		CID:         command.Sell,
		Callback:    m.Sell,
		Description: "Market",
	}
}

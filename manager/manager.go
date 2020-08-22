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
				Text:       "Please create a Adventurer by running the `create profile: YOURNAME` (no funny character in it => no space, no ponctuation)",
				ActionFlag: map[command.ID]string{},
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
					ActionFlag: map[command.ID]string{
						command.Profile: "See my profile",
					},
				},
				TownID: t.ID,
				Callback: map[command.ID]communication.ActionCallback{
					command.Profile: m.Profile,
				},
			}
		}
		createProfile := message.Command.(command.CommandCreateName)
		t.CreateAdventurer(createProfile.Name, message.PlayerID)
		result.Content = communication.ContentMessage{
			Text: "Your character is created! Welcome " + createProfile.Name,
			ActionFlag: map[command.ID]string{
				command.Profile: "See my profile",
			},
		}
		result.Callback = map[command.ID]communication.ActionCallback{
			command.Profile: m.Profile,
		}

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

func (m *Manager) ExploreDeep(message communication.ActionToManager, rid town.RegionID) *communication.ActionFromManager {
	a := m.Towns[message.TownID].Adventurers[message.PlayerID]

	for _, region := range m.Towns[message.TownID].Regions {
		if region.Name == rid {
			actions := GetAvailableActions(region.Name, region.Level, a, m.Towns[message.TownID], map[command.ID]communication.DescriptionAction{
				command.Profile: {
					CID:         command.Profile,
					Callback:    m.Profile,
					Description: "To profile",
				},
				command.Back: {
					CID:         command.Back,
					Callback:    m.Explore,
					Description: "Back",
				},
			})

			callback := map[command.ID]communication.ActionCallback{}
			actionsFlag := map[command.ID]string{}
			for _, action := range actions {
				callback[action.CID] = action.Callback
				actionsFlag[action.CID] = action.Description
			}

			return &communication.ActionFromManager{
				TownID: message.TownID,
				AllowList: []*player.ID{
					&message.PlayerID,
				},
				Callback: callback,
				Content: communication.ContentMessage{
					Text:       "You are exploring " + string(rid),
					ActionFlag: actionsFlag,
				},
				Parent: message.ParentID,
			}
		}
	}
	return &communication.ActionFromManager{
		TownID: message.TownID,
		AllowList: []*player.ID{
			&message.PlayerID,
		},
		Callback: map[command.ID]communication.ActionCallback{
			command.Profile: m.Profile,
			command.Explore: m.Explore,
		},
		Content: communication.ContentMessage{
			Text: "No action found in exploring " + string(rid),
			ActionFlag: map[command.ID]string{
				command.Profile: "Go back to the profile",
				command.Explore: "Go back to the exploration",
			},
		},
		Parent: message.ParentID,
	}

}

func (m *Manager) Explore(message communication.ActionToManager) *communication.ActionFromManager {

	callback := map[command.ID]communication.ActionCallback{}
	actionsFlag := map[command.ID]string{}
	for _, region := range m.Towns[message.TownID].Regions {
		callback[region.Command] = func(message communication.ActionToManager) *communication.ActionFromManager {
			return m.ExploreDeep(message, region.Name)
		}
		actionsFlag[region.Command] = "Explore in " + string(region.Name)
	}
	result := &communication.ActionFromManager{
		TownID: message.TownID,
		AllowList: []*player.ID{
			&message.PlayerID,
		},
		Callback: callback,
		Content: communication.ContentMessage{
			Text:       "Exploration asked.",
			ActionFlag: actionsFlag,
		},
		Parent: message.ParentID,
	}

	return result
}

func (m *Manager) Profile(message communication.ActionToManager) *communication.ActionFromManager {
	return &communication.ActionFromManager{
		TownID: message.TownID,
		AllowList: []*player.ID{
			&message.PlayerID,
		},
		Callback: map[command.ID]communication.ActionCallback{
			command.Explore: m.Explore,
			command.Craft:   m.Craft,
			command.Sell:    m.Sell,
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
		Callback: map[command.ID]communication.ActionCallback{
			// command.Fight: func(action communication.ActionToManager) *communication.ActionFromManager {
			// 	return m.CraftElement(action, town.WeaponCrafter)
			// },
			command.Profile: m.Profile,
		},
		Content: communication.ContentMessage{
			Text: "Craft (WIP)",
			ActionFlag: map[command.ID]string{
				//command.Fight:   "Ask the weapon crafter to craft a weapon for you",
				command.Profile: "Go back to your profile",
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
		Callback: map[command.ID]communication.ActionCallback{
			command.Explore: m.Explore,
			command.Craft:   m.Craft,
			command.Sell:    m.Sell,
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

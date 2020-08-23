package manager

import (
	"github.com/Selsynn/cbepbackend/business/characters"
	"github.com/Selsynn/cbepbackend/business/command"
	"github.com/Selsynn/cbepbackend/business/item"
	"github.com/Selsynn/cbepbackend/business/player"
	"github.com/Selsynn/cbepbackend/business/town"
	"github.com/Selsynn/cbepbackend/communication"
)

type RegionActionRegistry struct {
	ID             town.RegionID
	CID            command.ID
	Description    string
	ResultCallback string
	getRessources  func(lvl town.RegionLevel) []*item.Resources
	sideEffects    func(t *town.Town)
}

var regionList map[town.RegionID]map[town.RegionLevel][]*RegionActionRegistry

func init() {
	forest1 := []*RegionActionRegistry{
		{
			ID: town.Forest,
			getRessources: func(lvl town.RegionLevel) []*item.Resources {
				return []*item.Resources{
					{
						Item: item.Item{
							Name: item.Wood,
						},
						Qty: int(lvl),
					},
					{
						Item: item.Item{
							Name: item.Leather,
						},
						Qty: int(lvl) / 2,
					},
				}
			},
			CID:            command.Explore,
			Description:    "Explore the forest!",
			ResultCallback: "You have explored everything, found some wood and scraps of leather",
		},
		{
			ID: town.Forest,
			getRessources: func(lvl town.RegionLevel) []*item.Resources {
				return []*item.Resources{
					{
						Item: item.Item{
							Name: item.Leather,
						},
						Qty: int(lvl) * 2,
					},
				}
			},
			CID:            command.Protect,
			Description:    "Protect the people working in the forest from monsters",
			ResultCallback: "They were thankful of you and gave you some leather in thanks",
		},
	}
	forest5 := &RegionActionRegistry{
		ID: town.Forest,
		getRessources: func(lvl town.RegionLevel) []*item.Resources {
			return []*item.Resources{
				{
					Item: item.Item{
						Name: item.Wood,
					},
					Qty: int(lvl) * 10,
				},
			}
		},
		CID:            command.Build,
		Description:    "Exploit the forest",
		ResultCallback: "Working at the exploitation brings you a ton of wood",
		sideEffects: func(t *town.Town) {
			t.DeltaNature(-5)
		},
	}
	forest20 := &RegionActionRegistry{
		ID: town.Forest,
		getRessources: func(lvl town.RegionLevel) []*item.Resources {
			return []*item.Resources{}
		},
		CID:            command.Wood,
		Description:    "Dwelve deeper in the forest",
		ResultCallback: "You stumble into an enchanted forest",
		sideEffects: func(t *town.Town) {
			t.NewRegionLevel(town.EnchantedForest, command.EnchantedForest, 1)
		},
	}
	regionList = map[town.RegionID]map[town.RegionLevel][]*RegionActionRegistry{
		town.Forest: {
			town.RegionLevel(1):  forest1,
			town.RegionLevel(5):  append(forest1, forest5),
			town.RegionLevel(20): append(forest1, forest5, forest20),
		},
	}
}

func getContextCallback(rid town.RegionID,
	lvl town.RegionLevel,
	adventurer *characters.Adventurer,
	t *town.Town,
	contextActions map[command.ID]communication.DescriptionAction,
	rar RegionActionRegistry,
) communication.ActionCallback {

	return func(action communication.ActionToManager) *communication.ActionFromManager {
		adventurer.DeltaItems(rar.getRessources(lvl))
		if rar.sideEffects != nil {
			rar.sideEffects(t)
		}

		return &communication.ActionFromManager{
			Parent: action.ParentID,
			AllowList: []*player.ID{
				&action.PlayerID,
			},
			Callback: contextActions,
			Content: communication.ContentMessage{
				Text: rar.ResultCallback,
			},
			TownID: action.TownID,
		}
	}
}

func getAvailableActions(
	rid town.RegionID,
	lvl town.RegionLevel,
	adventurer *characters.Adventurer,
	t *town.Town,
	contextActions map[command.ID]communication.DescriptionAction,
) map[command.ID]communication.DescriptionAction {

	lvlCandidate := town.RegionLevel(-1)
	for unlockLvl := range regionList[rid] {
		if unlockLvl <= lvl && unlockLvl > lvlCandidate {
			lvlCandidate = unlockLvl
		}
	}
	if lvlCandidate == -1 {
		return map[command.ID]communication.DescriptionAction{}
	}
	ras := map[command.ID]communication.DescriptionAction{}
	for _, rar := range regionList[rid][lvlCandidate] {
		ras[rar.CID] = communication.DescriptionAction{
			CID:         rar.CID,
			Callback:    getContextCallback(rid, lvl, adventurer, t, contextActions, *rar),
			Description: rar.Description,
		}
	}
	return ras
}

func (m *Manager) exploreDeepContainer(regionID town.RegionID) communication.ActionCallback {
	return func(message communication.ActionToManager) *communication.ActionFromManager {
		return m.exploreDeep(message, regionID)
	}
}

func (m *Manager) Explore(message communication.ActionToManager) *communication.ActionFromManager {

	callback := map[command.ID]communication.DescriptionAction{}
	for _, region := range m.Towns[message.TownID].Regions {
		callback[region.Command] = communication.DescriptionAction{
			CID:         region.Command,
			Callback:    m.exploreDeepContainer(region.Name),
			Description: "Explore in " + string(region.Name),
		}
	}
	result := &communication.ActionFromManager{
		TownID: message.TownID,
		AllowList: []*player.ID{
			&message.PlayerID,
		},
		Callback: callback,
		Content: communication.ContentMessage{
			Text: "Exploration asked.",
		},
		Parent: message.ParentID,
	}

	return result
}

func (m *Manager) exploreDeep(message communication.ActionToManager, rid town.RegionID) *communication.ActionFromManager {
	a := m.Towns[message.TownID].Adventurers[message.PlayerID]

	for _, region := range m.Towns[message.TownID].Regions {
		if region.Name == rid {
			actions := getAvailableActions(region.Name, region.Level, a, m.Towns[message.TownID], map[command.ID]communication.DescriptionAction{
				command.Profile: m.ProfileCallback(),
				command.Back: {
					CID:         command.Back,
					Callback:    m.Explore,
					Description: "Back",
				},
			})

			return &communication.ActionFromManager{
				TownID: message.TownID,
				AllowList: []*player.ID{
					&message.PlayerID,
				},
				Callback: actions,
				Content: communication.ContentMessage{
					Text: "You are exploring " + string(rid),
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
		Callback: map[command.ID]communication.DescriptionAction{
			command.Profile: m.ProfileCallback(),
			command.Explore: m.ExploreCallback(),
		},
		Content: communication.ContentMessage{
			Text: "No action found in exploring " + string(rid),
		},
		Parent: message.ParentID,
	}

}

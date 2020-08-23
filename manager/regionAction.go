package manager

import (
	"github.com/Selsynn/cbepbackend/business/command"
	"github.com/Selsynn/cbepbackend/business/item"
	"github.com/Selsynn/cbepbackend/business/player"
	"github.com/Selsynn/cbepbackend/business/town"
	"github.com/Selsynn/cbepbackend/communication"
)

// type RegionAction struct {
// 	ID          town.RegionID
// 	CID         command.ID
// 	Description string
// 	Callback    communication.ActionCallback
// }

//func (ra *RegionAction) ThisIsARegionAction() {}

type RegionActionRegistry struct {
	ID             town.RegionID
	CID            command.ID
	Description    string
	ResultCallback string
	getRessources  func(lvl town.RegionLevel) []*town.Resources
	sideEffects    func(t *town.Town)
}

var regionList map[town.RegionID]map[town.RegionLevel][]*RegionActionRegistry

func init() {
	forest1 := []*RegionActionRegistry{
		{
			ID: town.Forest,
			getRessources: func(lvl town.RegionLevel) []*town.Resources {
				return []*town.Resources{
					{
						Item: town.Item{
							Name: item.Wood,
						},
						Qty: int(lvl),
					},
					{
						Item: town.Item{
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
			getRessources: func(lvl town.RegionLevel) []*town.Resources {
				return []*town.Resources{
					{
						Item: town.Item{
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
		getRessources: func(lvl town.RegionLevel) []*town.Resources {
			return []*town.Resources{
				{
					Item: town.Item{
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
			t.Nature--
		},
	}
	forest20 := &RegionActionRegistry{
		ID: town.Forest,
		getRessources: func(lvl town.RegionLevel) []*town.Resources {
			return []*town.Resources{}
		},
		CID:            command.Wood,
		Description:    "Dwelve deeper in the forest",
		ResultCallback: "You stumble into an enchanted forest",
		sideEffects: func(t *town.Town) {
			for _, r := range t.Regions {
				if r.Name == town.EnchantedForest {
					return
				}
			}
			t.Regions = append(t.Regions, &town.Region{
				Name:    town.EnchantedForest,
				Command: command.EnchantedForest,
				Level:   1,
			})
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

func GetContextCallback(rid town.RegionID,
	lvl town.RegionLevel,
	adventurer *town.Adventurer,
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

func GetAvailableActions(
	rid town.RegionID,
	lvl town.RegionLevel,
	adventurer *town.Adventurer,
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
			Callback:    GetContextCallback(rid, lvl, adventurer, t, contextActions, *rar),
			Description: rar.Description,
		}
	}
	return ras
}

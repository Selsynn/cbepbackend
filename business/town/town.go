package town

import (
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/Selsynn/cbepbackend/business/characters"
	"github.com/Selsynn/cbepbackend/business/command"
	"github.com/Selsynn/cbepbackend/business/item"
	"github.com/Selsynn/cbepbackend/business/player"
)

type ID string

type Town struct {
	ID                ID
	Name              string
	Villagers         int
	Threat            int
	LastThreatChecked time.Time
	Nature            int
	NPC               map[characters.Profession]*characters.NPC //	Merchants         []*NPC
	Upgrades          []*Upgrade
	Adventurers       map[player.ID]*characters.Adventurer
	Regions           []*Region
}

type (
	UpgradeName string
	Upgrade     struct {
		Name UpgradeName
	}
)

func New() *Town {
	t := &Town{
		NPC:               characters.InitNPCForNewCity(),
		Upgrades:          []*Upgrade{},
		Adventurers:       map[player.ID]*characters.Adventurer{},
		Name:              "IdleTown",
		LastThreatChecked: time.Now(),
		Nature:            0,
		Threat:            0,
		Villagers:         3,
		ID:                ID(uuid.New().String()),
		Regions: []*Region{
			{
				Name:    Forest,
				Command: command.Wood,
				Level:   RegionLevel(1),
			},
		},
	}

	return t
}

func (t *Town) CreateAdventurer(name string, playerID player.ID) {
	t.Adventurers[playerID] = &characters.Adventurer{
		Name:          name,
		ID:            (uuid.New().String()),
		Relationships: make(map[*characters.NPC]int),
		Items:         map[item.ID]*item.Resources{},
	}

	fmt.Printf("Create Adventurer \n%#v\n", t.Adventurers)
}

func (t *Town) DeltaNature(delta int) {
	t.Nature += delta
}

// func (t *Town) Craft(c *characters.NPC, i *item.Item, a *characters.Adventurer) {
// 	fmt.Printf("You asked for crafting %#v by crafter %s \n", i, c.Name)
// 	c.AddWork(i, 10, a)
// }

// func (t *Town) GetItem(name item.ID) *item.Item {
// 	switch name {
// 	case item.Bow:
// 		return &item.Item{
// 			Name: item.Bow,
// 		}
// 	}
// 	return nil
// }

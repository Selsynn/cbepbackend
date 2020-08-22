package town

import (
	"fmt"
	"time"

	"github.com/google/uuid"

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
	NPC               map[Profession]*NPC //	Merchants         []*NPC
	Upgrades          []*Upgrade
	Adventurers       map[player.ID]*Adventurer
	Regions           []*Region
}

type (
	Adventurer struct {
		Name          string
		Relationships map[*NPC]int
		ID            string
		Consummable   []*Resources
		Items         []*Resources
		Resources     []*Resources
	}

	Item struct {
		Name item.ID
	}

	Resources struct {
		Item
		Qty int
	}

	UpgradeName string
	Upgrade     struct {
		Name UpgradeName
	}

	RegionName string
	Region     struct {
		Name    string
		Level   int
		Command command.ID
		Actions map[command.ID]RegionAction
	}

	RegionAction interface {
		ThisIsARegionAction()
	}
)

func New() *Town {
	t := &Town{
		NPC: map[Profession]*NPC{
			Merchant: {
				lastWorkedTime: time.Now(),
				workDone:       make([]WorkItem, 0),
				workQueue:      []WorkItem{},
				Name:           "Gripsou",
				RelQuest:       map[int]Quest{},
				Specialty:      Merchant,
			},
			WeaponCrafter: {
				Name:           "Kreator",
				lastWorkedTime: time.Now(),
				workDone:       make([]WorkItem, 0),
				workQueue:      []WorkItem{},
				RelQuest:       map[int]Quest{},
				Specialty:      WeaponCrafter,
				Knowledge: []WorkItem{
					{
						Item: &Item{
							Name: item.Bow,
						},
						Cost: []Resources{
							{
								Item: Item{
									Name: item.Wood,
								},
								Qty: 10,
							},
						},
						WorkCost: 1,
					},
				},
			},
		},
		Upgrades:          []*Upgrade{},
		Adventurers:       map[player.ID]*Adventurer{},
		Name:              "IdleTown",
		LastThreatChecked: time.Now(),
		Nature:            0,
		Threat:            0,
		Villagers:         3,
		ID:                ID(uuid.New().String()),
		Regions:           []*Region{},
	}

	//t.CreateMerchant()
	//t.CreateCrafter()
	//t.CreateAdventurer()

	return t
}

// func (t Town) DescribeCity() string {
// 	fmt.Printf("Describe CITY \n%#v\n", t)

// 	return fmt.Sprintf("%s has %d brave adventurers, %d crafters/ merchants.", t.Name, len(t.Adventurers), len(t.NPC))
// }

// func (t Town) DescribeHeros() string {
// 	return fmt.Sprintf("%s has %d brave adventurers, %d crafters/ merchants.", t.Name, len(t.Adventurers), len(t.NPC))
// }
// func (t *Town) CreateMerchant() {
// 	t.NPC = append(t.NPC, &NPC{
// 		Name:      "Super Generic Merchant",
// 		RelQuest:  make(map[int]Quest),
// 		Specialty: Merchant,
// 	})

// 	fmt.Printf("Create MERCHANT \n%#v\n", t.Merchants)
// }

// func (t *Town) CreateCrafter() {
// 	t.Crafters = append(t.Crafters, &NPC{
// 		Name:      "Super Generic Crafter",
// 		RelQuest:  make(map[int]Quest),
// 		Specialty: WeaponCrafter,
// 	})

// 	fmt.Printf("Create Crafter \n%#v\n", t.Crafters)
// }

func (t *Town) CreateAdventurer(name string, playerID player.ID) {
	t.Adventurers[playerID] = &Adventurer{
		Name:          name,
		ID:            (uuid.New().String()),
		Relationships: make(map[*NPC]int),
		Consummable:   []*Resources{},
		Items:         []*Resources{},
		Resources:     []*Resources{},
	}

	fmt.Printf("Create Adventurer \n%#v\n", t.Adventurers)
}

func (t *Town) Craft(c *NPC, i *Item, a *Adventurer) {
	fmt.Printf("You asked for crafting %#v by crafter %s \n", i, c.Name)
	c.AddWork(i, 10, a)
}

func (t *Town) GetItem(name item.ID) *Item {
	switch name {
	case item.Bow:
		return &Item{
			Name: item.Bow,
		}
	}
	return nil
}

func (r *Resources) Describe() string {
	return fmt.Sprintf("(%d *%s*)", r.Qty, r.Name)
}

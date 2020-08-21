package town

import (
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/Selsynn/DiscordBotTest1/business/item"
)

type ID string

type Town struct {
	ID                ID
	Name              string
	Villagers         int
	Threat            int
	LastThreatChecked time.Time
	Nature            int
	Crafters          []*NPC
	Merchants         []*NPC
	Consummable       []*Item
	Items             []*Item
	Resources         []*Item
	Upgrades          []*Upgrade
	Adventurers       []*Adventurer
}

type (
	Adventurer struct {
		Name          string
		Relationships map[*NPC]int
		ID            string
	}

	Item struct {
		Name item.ID
	}

	UpgradeName string
	Upgrade     struct {
		Name UpgradeName
	}
)

func New() *Town {
	t := &Town{
		Crafters:          []*NPC{},
		Merchants:         []*NPC{},
		Consummable:       []*Item{},
		Items:             []*Item{},
		Resources:         []*Item{},
		Upgrades:          []*Upgrade{},
		Adventurers:       []*Adventurer{},
		Name:              "IdleTown",
		LastThreatChecked: time.Now(),
		Nature:            0,
		Threat:            0,
		Villagers:         3,
		ID:                ID(uuid.New().String()),
	}

	t.CreateMerchant()
	t.CreateCrafter()
	t.CreateAdventurer()

	return t
}

func (t Town) DescribeCity() string {
	fmt.Printf("Describe CITY \n%#v\n", t)

	return fmt.Sprintf("%s has %d brave adventurers, %d crafters, %d merchants.", t.Name, len(t.Adventurers), len(t.Crafters), len(t.Merchants))
}

func (t Town) DescribeHeros() string {
	return fmt.Sprintf("%s has %d brave adventurers, %d crafters, %d merchants.", t.Name, len(t.Adventurers), len(t.Crafters), len(t.Merchants))
}

func (t *Town) CreateMerchant() {
	t.Merchants = append(t.Merchants, &NPC{
		Name:      "Super Generic Merchant",
		RelQuest:  make(map[int]Quest),
		Specialty: Merchant,
	})

	fmt.Printf("Create MERCHANT \n%#v\n", t.Merchants)
}

func (t *Town) CreateCrafter() {
	t.Crafters = append(t.Crafters, &NPC{
		Name:      "Super Generic Crafter",
		RelQuest:  make(map[int]Quest),
		Specialty: Crafter,
	})

	fmt.Printf("Create Crafter \n%#v\n", t.Crafters)
}

func (t *Town) CreateAdventurer() {
	t.Adventurers = append(t.Adventurers, &Adventurer{
		Name:          "SuperAdventurer",
		ID:            "SomeSecretId",
		Relationships: make(map[*NPC]int),
	})

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

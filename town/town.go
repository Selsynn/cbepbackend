package town

import (
	"fmt"
)

type Profession string

const (
	Merchant Profession = "Merchant"
	Crafter  Profession = "Crafter"
)

type Town struct {
	Name              string
	Villagers         int
	Threat            int
	LastThreatChecked int //timestamp
	Nature            int
	Crafters          []NPC
	Merchants         []NPC
	Consummable       []Item
	Items             []Item
	Resources         []Item
	Upgrades          []Upgrade
	Adventurers       []Adventurer
}

type (
	NPCName string
	NPC     struct {
		Name      NPCName
		Specialty Profession
		RelQuest  map[int]Quest
	}
	Adventurer struct {
		Name          string
		Relationships map[*NPC]int
	}

	ItemName string
	Item     struct {
		Name ItemName
	}

	UpgradeName string
	Upgrade     struct {
		Name UpgradeName
	}

	Quest struct {
	}
)

func New() *Town {
	t := &Town{
		Crafters:    []NPC{},
		Merchants:   []NPC{},
		Consummable: []Item{},
		Items:       []Item{},
		Resources:   []Item{},
		Upgrades:    []Upgrade{},
		Adventurers: []Adventurer{},
	}
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
	t.Merchants = append(t.Merchants, NPC{
		Name:      "Super Generic Merchant",
		RelQuest:  make(map[int]Quest),
		Specialty: Merchant,
	})

	fmt.Printf("Create MERCHANT \n%#v\n", t.Merchants)
}

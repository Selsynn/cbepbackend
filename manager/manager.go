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
		t = town.New()
		m.Towns[message.TownID] = t
		fmt.Printf("Creation of a new city \n%#v\n", t)
	}

	a := t.Adventurers[0]

	fmt.Printf("Describe WORLD \n%#v\n", m.Towns)

	result := &communication.ActionFromManager{
		TownID: message.TownID,
		AllowList: []*player.ID{
			&message.PlayerID,
		},
	}

	var resultContent string

	switch message.Command.ID() {
	case command.ViewShop:
		resultContent = t.DescribeCity()
	case command.NewMerchant:
		t.CreateMerchant()
		resultContent = "Merchant created." + t.DescribeCity()
		result.Callback = map[command.ID]func() *communication.ActionFromManager{
			command.Accept: func() *communication.ActionFromManager {
				fmt.Println("Accepted !")
				return nil
			},
		}
	case command.ViewHeros:
		resultContent = "You asked to see the heros" + t.DescribeHeros()
	case command.Craft:
		craft := message.Command.(command.CommandCraft)
		t.Craft(t.Crafters[0], t.GetItem(craft.ItemID), a)
	default:
		resultContent = fmt.Sprintf("Il n'y a rien a cette adresse. List of all the command currently supported: **%v**", command.ListAll())
	}

	result.Content = communication.ContentMessage{
		Text: resultContent,
	}

	return result
}

func NewManager() Manager {
	return Manager{
		Towns: make(map[town.ID]*town.Town),
	}
}

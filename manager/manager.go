package manager

import (
	"fmt"

	"github.com/Selsynn/cbepbackend/business/command"
	businessTown "github.com/Selsynn/cbepbackend/business/town"
	"github.com/Selsynn/cbepbackend/communication"
	"github.com/Selsynn/cbepbackend/town"
)

type Manager struct {
	Towns map[businessTown.ID]*town.Town
}

func (m *Manager) Process(message communication.ActionToManager) communication.ActionFromManager {
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

	result := communication.ActionFromManager{
		TownID: message.TownID,
	}

	var resultContent string

	switch message.Command.ID() {
	case command.ViewShop:
		resultContent = t.DescribeCity()
	case command.NewMerchant:
		t.CreateMerchant()
		resultContent = "Merchant created." + t.DescribeCity()
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

	// index := strings.Index(message.Content, ":")
	// command := message.Content
	// if index != -1 {
	// 	command = message.Content[:index]
	// }

	// switch command {
	// case command_shop:
	// 	result.Content = t.DescribeCity()
	// case create_merchant:
	// 	t.CreateMerchant()
	// 	result.Content = "Merchant created." + t.DescribeCity()
	// case command_heros:
	// 	result.Content = "You asked to see the heros" + t.DescribeHeros()
	// case command_craft:
	// 	t.Craft(t.Crafters[0], t.GetItem(message.Content[(index+1):]), a)
	// default:
	// 	result.Content = "Il n'y a rien a cette adresse. List of all the command currently supported: **" + command_shop + "**, **" + command_heros + "**"
	// }

	return result
}

func NewManager() Manager {
	return Manager{
		Towns: make(map[businessTown.ID]*town.Town),
	}
}

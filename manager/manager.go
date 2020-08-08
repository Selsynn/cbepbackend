package manager

import (
	"fmt"

	"github.com/Selsynn/DiscordBotTest1/talker"
	"github.com/Selsynn/DiscordBotTest1/town"
)

type Manager struct {
	Towns map[talker.Server]*town.Town
}

const command_shop = "shop"
const command_heros = "heros"
const create_merchant = "new merchant"

func (m *Manager) Process(message talker.Message) talker.Order {
	//Get town
	t, found := m.Towns[message.Server]

	if !found {
		//No city found
		t = town.New()
		m.Towns[message.Server] = t
		fmt.Printf("Creation of a new city \n%#v\n", t)
	}

	fmt.Printf("Describe WORLD \n%#v\n", m.Towns)

	result := talker.Order{
		Write: message.Write,
	}

	switch message.Content {
	case command_shop:
		result.Content = t.DescribeCity()
	case create_merchant:
		t.CreateMerchant()
		result.Content = "Merchant created." + t.DescribeCity()
	case command_heros:
		result.Content = "You asked to see the heros" + t.DescribeHeros()
	default:
		result.Content = "Il n'y a rien a cette adresse. List of all the command currently supported: **" + command_shop + "**, **" + command_heros + "**"
	}

	return result
}

func NewManager() Manager {
	return Manager{
		Towns: make(map[talker.Server]*town.Town),
	}
}

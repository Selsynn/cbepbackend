package manager

import (
	"fmt"
	"strings"

	"github.com/Selsynn/DiscordBotTest1/talker"
	"github.com/Selsynn/DiscordBotTest1/town"
)

type Manager struct {
	Towns map[talker.Server]*town.Town
}

const command_shop = "shop"
const command_heros = "heros"
const create_merchant = "new merchant"
const command_craft = "craft" //craft:bow

func (m *Manager) Process(message talker.Message) talker.Order {
	//Get town
	t, found := m.Towns[message.Server]

	if !found {
		//No city found
		t = town.New()
		m.Towns[message.Server] = t
		fmt.Printf("Creation of a new city \n%#v\n", t)
	}

	a := t.Adventurers[0]

	fmt.Printf("Describe WORLD \n%#v\n", m.Towns)

	result := talker.Order{
		Write: message.Write,
	}

	index := strings.Index(message.Content, ":")
	command := message.Content
	if index != -1 {
		command = message.Content[:index]
	}

	switch command {
	case command_shop:
		result.Content = t.DescribeCity()
	case create_merchant:
		t.CreateMerchant()
		result.Content = "Merchant created." + t.DescribeCity()
	case command_heros:
		result.Content = "You asked to see the heros" + t.DescribeHeros()
	case command_craft:
		t.Craft(t.Crafters[0],t.GetItem(message.Content[(index+1):]), a)
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

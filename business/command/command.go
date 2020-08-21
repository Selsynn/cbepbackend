package command

import (
	"fmt"
	"strings"

	"github.com/Selsynn/cbepbackend/business/item"
)

type ID string

const (
	//BASIC CMD
	Explore ID = "Explore"
	Profile ID = "Profile"
	//CONTEXT CMD

	// ViewShop    ID = "shop"
	// ViewHeros   ID = "heros"
	// NewMerchant ID = "new merchant"
	Craft   ID = "Craft"
	Accept  ID = "Accept"
	Refuse  ID = "Refuse"
	Fight   ID = "Fight"
	Build   ID = "Build"
	Protect ID = "Protect"
	Sell    ID = "Sell"
	Wood    ID = "Wood"
)

type Command interface {
	ID() ID
}

type CommandSimple struct {
	Id ID
}

func (c CommandSimple) ID() ID {
	return c.Id
}

type CommandCraft struct {
	CommandSimple
	ItemID item.ID
}

func ListAll() []ID {
	return []ID{
		// ViewHeros,
		// ViewShop,
		// NewMerchant,
		Craft,
	}
}

func Parse(text string) (Command, error) {
	index := strings.Index(text, ":")
	t := text
	if index != -1 {
		t = text[:index]
	}
	t = strings.TrimSpace(t)
	t = strings.ToUpper(t[:1]) + t[1:]
	id := ID(t)
	switch id {
	case Explore, Profile:
		return CommandSimple{
			Id: id,
		}, nil

	// case ViewHeros,
	// 	ViewShop,
	// 	NewMerchant:
	// return CommandSimple{
	// 	Id: id,
	// }, nil
	case Craft:
		return CommandCraft{
			CommandSimple: CommandSimple{
				Id: id,
			},
			ItemID: item.ID(text[(index + 1):]),
		}, nil
	default:
		return nil, fmt.Errorf("Command not found for %s\n", text)
	}
}

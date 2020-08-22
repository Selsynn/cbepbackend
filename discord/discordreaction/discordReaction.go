package discordreaction

import "github.com/Selsynn/cbepbackend/business/command"

type ID string

// const (
// 	//When adding a emoji, you have to also add it in the command list
// 	   "✔":Accept,
// 	   "❌":Refuse,
// 	    "⚔":Fight,
// 	    "⚒":Build,
// 	  "🛡":Protect,
// 	     "⚖":Sell,
// )

// func GetAll() []ID {
// 	return []ID{
// 		Accept,
// 		Refuse,
// 		Fight,
// 		Build,
// 		Protect,
// 		Sell,
// 	}
// }

var reaction2Command = map[ID]command.ID{}
var command2Reaction = map[command.ID]ID{}

func init() {
	reaction2Command = map[ID]command.ID{
		"☑": command.Accept,
		"❌": command.Refuse,
		"⚔": command.Fight,
		"⚒": command.Build,
		"🛡": command.Protect,
		"⚖": command.Sell,
		"🌳": command.Wood,
		"🧑": command.Profile,
		"👀": command.Explore,
		"⚙": command.Craft,
		"🏹": command.Bow,
		"🔙": command.Back,
		"🌇": command.Town,
		"🌲": command.EnchantedForest,
	}

	command2Reaction = map[command.ID]ID{}

	for reaction, command := range reaction2Command {
		command2Reaction[command] = reaction
	}
}

func Match2Command(id ID) command.ID {
	return reaction2Command[id]
}

func Match2Reaction(id command.ID) ID {
	return command2Reaction[id]
}

package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/Selsynn/cbepbackend/discord"
	"github.com/Selsynn/cbepbackend/interaction"
	"github.com/Selsynn/cbepbackend/interaction/interactiondiscordimpl"
	"github.com/Selsynn/cbepbackend/manager"
	"github.com/Selsynn/cbepbackend/talker"
	"github.com/Selsynn/cbepbackend/talker/talkercmd"
	"github.com/Selsynn/cbepbackend/talker/talkerdiscord"
)

// Variables used for command line parameters
var (
	Token             string
	BotImplementation string
)

func initMain() {

	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.StringVar(&BotImplementation, "i", "discord", "Bot Implementation['discord','cmd']")
	flag.Parse()
}

func main() {
	initMain()

	var t talker.Talker
	var i interaction.Interaction

	switch BotImplementation {
	case "discord":
		fmt.Printf("BotImplementation chosen: Discord\n")
		var shutdownTalker func()
		t, shutdownTalker = talkerdiscord.NewTalkerDiscord(Token)
		i = &interactiondiscordimpl.InteractionDiscord{
			Servers: make(map[discord.ServerID]discord.Server),
		}
		defer shutdownTalker()
	case "cmd":
		fmt.Printf("BotImplementation chosen: Cmd\n")
		t = talkercmd.New()
	default:
		panic("verify parameters")
	}

	m := manager.NewManager()
	go func() {
		for mess := range t.Read() {
			toManager := i.GetActionToManager(mess)
			if toManager.Command == nil {
				continue
			}
			processed := m.Process(toManager)
			fromManager := i.GetActionFromManager(processed)
			t.Write(fromManager)
		}
		// 	order := m.Process(mess)
		// 	order.Write(order.Content)
		// }
	}()

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

}

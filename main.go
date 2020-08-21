package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/Selsynn/craft-build-explore-protect-backend/discord"
	"github.com/Selsynn/craft-build-explore-protect-backend/interaction"
	"github.com/Selsynn/craft-build-explore-protect-backend/interaction/interactiondiscordimpl"
	"github.com/Selsynn/craft-build-explore-protect-backend/manager"
	"github.com/Selsynn/craft-build-explore-protect-backend/talker"
	"github.com/Selsynn/craft-build-explore-protect-backend/talker/talkercmd"
	"github.com/Selsynn/craft-build-explore-protect-backend/talker/talkerdiscord"
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
			t.Write(i.GetActionFromManager(m.Process(i.GetActionToManager(mess))))
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

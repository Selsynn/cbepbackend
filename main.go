package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/Selsynn/DiscordBotTest1/manager"
	"github.com/Selsynn/DiscordBotTest1/talker"
	"github.com/Selsynn/DiscordBotTest1/talker/talkercmd"
	"github.com/Selsynn/DiscordBotTest1/talker/talkerdiscord"
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

	switch BotImplementation {
	case "discord":
		var shutdownTalker func()
		t, shutdownTalker = talkerdiscord.NewTalkerDiscord(Token)
		defer shutdownTalker()
	case "cmd":
		t = talkercmd.New()
	default:
		panic("verify parameters")
	}

	m := manager.NewManager()
	go func() {
		for mess := range t.Read() {
			order := m.Process(mess)
			order.Write(order.Content)
		}
	}()

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

}

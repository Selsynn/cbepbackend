package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/Selsynn/DiscordBotTest1/manager"

	"github.com/Selsynn/DiscordBotTest1/talker"
)

// Variables used for command line parameters
var (
	Token string
)

func initMain() {

	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.Parse()
}

func main() {
	initMain()

	t, shutdownTalker := talker.NewTalkerDiscord(Token)
	defer shutdownTalker()

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
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

}

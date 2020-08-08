package main

import (
	"testing"
)

func TestIntegration(t *testing.T) {
	// initMain()
	// // Ignore all messages created by the bot itself
	// // This isn't required in this specific example but it's a good practice.
	// // if m.Author.ID == s.State.User.ID {
	// // 	return
	// // }
	// s := &discordgo.Session{
	// 	State: &discordgo.State{
	// 		Ready: discordgo.Ready{
	// 			User: &discordgo.User{
	// 				ID: "UserId",
	// 			},
	// 		},
	// 	},
	// }
	// m := &discordgo.MessageCreate{
	// 	Message: &discordgo.Message{
	// 		Author: &discordgo.User{
	// 			ID: "Bot",
	// 		},
	// 		Content: "pingoiu",
	// 	},
	// }
	// messageCreate(s, m)

	// got := Abs(-1)
	// if got != 1 {
	// 	t.Errorf("Abs(-1) = %d; want 1", got)
	// }
}

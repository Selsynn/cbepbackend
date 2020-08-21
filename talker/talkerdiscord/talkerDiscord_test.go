package talkerdiscord

import (
	"os"
	"testing"
	"time"

	"github.com/Selsynn/craft-build-explore-protect-backend /talker"
	"github.com/bwmarrin/discordgo"
)

// func TestWrite(tester *testing.T) {
// 	t := TalkerDiscord{}
// 	o := Order{}

// 	t.Write(o)
// }

// func TestRead(tester *testing.T) {
// 	t := TalkerDiscord{}
// 	result := <-t.Read()

// 	fmt.Print("result", result)
// }

func testNew() talker.Talker {
	t, _ := NewTalkerDiscord("")
	return t
}

// func testNewServer() talker.Server {
// 	s := ServerDiscord{
// 		channelId: "",
// 		name:      "Name",
// 	}
// 	return s
// }

func Shup() {
	testNew()
	//testNewServer()
}

func TestWriteFormat(tester *testing.T) {
	key := os.Getenv("DiscordBotToken")

	if key == "" {
		tester.Skip("No key provided for Bot Discord - Skipping")
	}

	t, _ := NewTalkerDiscord(key)

	e := &discordgo.MessageEmbed{
		//Author: "SuperShop",
		Color:       5676621,
		Timestamp:   time.Now().UTC().Format("2006-01-02T15:04:05Z"),
		Description: "C'est une des-des-description!",
		Fields: []*discordgo.MessageEmbedField{
			{
				Inline: true,
				Name:   "SuperName",
				Value:  "Super value trop bien",
			},
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Ceci est un footer " + time.Now().Format("15:04:05") + "\nbliblou",
		},
		Image: &discordgo.MessageEmbedImage{
			Height: 10,
			URL:    "https://www.linuxmint.com/img/ads/info.png",
		},
	}

	_, err := t.session.ChannelMessageSendEmbed("735167412869005313", e)
	if err == nil {
		return
	}
	panic(err.Error())

}

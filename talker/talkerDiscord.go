package talker

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

type TalkerDiscord struct {
	session   *discordgo.Session
	messageCh chan Message
	Servers   []*ServerDiscord
}

type ServerDiscord struct {
	channelId string
	name      string
}

func NewTalkerDiscord(token string) (impl *TalkerDiscord, shutdown func()) {
	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		panic("error creating Discord session," + err.Error())
	}

	t := &TalkerDiscord{
		session:   dg,
		messageCh: make(chan Message, 1),
		Servers:   make([]*ServerDiscord, 0),
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate(t))

	// Register the ReactionGet all
 	dg.AddHandler(messageReactionAdd(t))

	dg.AddHandler(func(s *discordgo.Session, m *discordgo.MessageReactionRemove) {
		fmt.Printf("Remove %#v\n", m.MessageReaction)
	})
	dg.AddHandler(func(s *discordgo.Session, m *discordgo.MessageReactionRemoveAll) {
		fmt.Printf("Remove all %#v\n", m.MessageReaction.Emoji)
	})

	// In this example, we only care about receiving message events.
	dg.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuildMessages | discordgo.IntentsGuildMessageReactions)

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		panic("error opening connection," + err.Error())
	}

	return t, t.Shutdown
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func messageCreate(t *TalkerDiscord) func(s *discordgo.Session, m *discordgo.MessageCreate) {
	return func(s *discordgo.Session, m *discordgo.MessageCreate) {
		// Ignore all messages created by the bot itself
		// This isn't required in this specific example but it's a good practice.
		if m.Author.ID == s.State.User.ID {
			return
		}

		write := func(content string) {
			message, err := s.ChannelMessageSend(m.ChannelID, content)

			if err != nil {
				panic(err.Error())
			}

			err = s.MessageReactionAdd(m.ChannelID, message.ID, "🧪")

			if err != nil {
				fmt.Println(err.Error())
			}
		}

		mess := Message{
			Write:   write,
			Content: m.Content,
			Server:  t.FindOrCreateServer(m.ChannelID),
		}

		t.messageCh <- mess
	}

}

func messageReactionAdd(t *TalkerDiscord) func(s *discordgo.Session, m *discordgo.MessageReactionAdd) {
	return func(s *discordgo.Session, m *discordgo.MessageReactionAdd) {
		fmt.Printf("Add %#v\n", m.MessageReaction)
	}
}

func (t TalkerDiscord) Read() chan Message {
	return t.messageCh
}

func (t TalkerDiscord) Write(o Order) {
	o.Write(o.Content)
}

func (t TalkerDiscord) Shutdown() {
	// Cleanly close down the Discord session.
	t.session.Close()
}

func (t *TalkerDiscord) FindOrCreateServer(channelId string) *ServerDiscord {
	for _, item := range t.Servers {
		if item.GetId() == channelId {
			return item
		}
	}

	//we didn't find anything, time to create it
	s := &ServerDiscord{
		channelId: channelId,
		name:      "Basic Name",
	}
	t.Servers = append(t.Servers, s)
	return s
}

// func (t TalkerDiscord) GetServers() []Server {
// 	return t.Servers
// }

func (s ServerDiscord) GetName() string {
	return s.name
}

func (s ServerDiscord) GetId() string {
	return s.channelId
}

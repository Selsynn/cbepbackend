package talkercmd

import (
	"fmt"
	"time"

	"github.com/Selsynn/DiscordBotTest1/talker"
)

type TalkerCmd struct {
	messagesCh                      chan talker.Message
	discussionExpectedFromManagerCh chan string
	discussionSentToManagerCh       chan string
}

func New() *TalkerCmd {
	t := &TalkerCmd{
		messagesCh:                      make(chan talker.Message, 100),
		discussionExpectedFromManagerCh: make(chan string, 100),
		discussionSentToManagerCh:       make(chan string, 100),
	}

	t.discussionSentToManagerCh <- "This is the beginning! Do you want to continue (type yes)"
	t.discussionExpectedFromManagerCh <- "Il n'y a rien a cette adresse. List of all the command currently supported: **shop**, **heros**"
	t.discussionSentToManagerCh <- "shop"
	t.discussionExpectedFromManagerCh <- "IdleTown has 1 brave adventurers, 1 crafters, 1 merchants."
	t.discussionSentToManagerCh <- "craft:BOW"
	t.discussionExpectedFromManagerCh <- "IdleTown has 1 brave adventurers, 1 crafters, 1 merchants."
	close(t.discussionSentToManagerCh)
	close(t.discussionExpectedFromManagerCh)

	go t.autoAnswerToManager()
	return t
}

func (t TalkerCmd) autoAnswerToManager() {
	for sent := range t.discussionSentToManagerCh {
		time.Sleep(time.Second)
		fmt.Printf("TalkerCmd - Reading <<%s>> \n", sent)
		t.messagesCh <- talker.Message{
			Content: sent,
			Write: func(s string) {
				expected, open := <-t.discussionExpectedFromManagerCh
				if !open {
					fmt.Println("TalkerCmd - End of discussionExpectedFromManagerCh")
					panic("nothing expected")
				}
				if expected != s {
					panic(fmt.Sprintf("\nExpected \t<<%s>>\ngot \t\t<<%s>>", expected, s))
				}
				t.Write(talker.Order{
					Content: s,
				})
			},
		}
	}
	fmt.Println("TalkerCmd - End of discussionSentToManagerCh")
}

func (t TalkerCmd) Read() chan talker.Message {
	return t.messagesCh
}

func (t TalkerCmd) Write(o talker.Order) {
	fmt.Printf("TalkerCmd - Sending <<%s>> \n", o.Content)
}

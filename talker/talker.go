package talker

import "github.com/Selsynn/cbepbackend/communication"

type Talker interface {
	Read() chan MessageReceived
	Write(MessageSent) communication.ActionID
}

// type Message struct {
// 	Content string
// 	Write   func(string)
// 	Server  Server
// }

// type Order struct {
// 	Write   func(string)
// 	Content string
// }

type MessageReceived interface {
	ThisStructIsAMessageReceive()
}
type MessageSent interface {
	ThisStructIsAMessageSent()
}

package talker

type Talker interface {
	Read() chan MessageReceived
	Write(MessageSent)
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

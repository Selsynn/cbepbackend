package talker

type Talker interface {
	Read() chan Message
	Write(o Order)
	//GetServers() []Server
}

type Message struct {
	Content string
	Write   func(string)
	Server  Server
}

type Order struct {
	Write   func(string)
	Content string
}

type Server interface {
	GetName() string
	GetId() string
}

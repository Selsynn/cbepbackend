package interaction

type Action struct {
	Parent    *Action
	Content   ContentMessage
	AllowList []*UserID
	Callback  []func(Action)
	//CleanUp  func()
}

type ContentMessage interface {
	ToString() string
	AddAction(a Action)
	GetAllActions() []*Action
}

type UserID string

package simplechat

import "container/list"

type Message struct {
	ID   string
	From string
	To   string
	Type string // user group
	Data string
}

type connection struct {
	// Connection
}

type connectionContext struct {
	ID   string
	Type int
}

var (
	// _conn  map[string]connection
	_user  map[string]chatUser
	_group map[string]chatGroup
)

func init() {
	_user = make(map[string]chatUser)
	_group = make(map[string]chatGroup)
}

type chatUser struct {
	ID   string
	conn connection
}

type chatGroup struct {
	ID    string
	Admin string
	conn  connection

	userList *list.Element
}

type Client interface {
	NewClient(id string)
	ClientID() string
	Client() error
	SendToUser() error
	SendToGroup() error
	Role() int
}

type Group interface {
	NewGroup(id string)
	GroupID() string
	AddUserToGroup() error
	DeleteUser() error
	GetUserList() error
	SendMessage() error
}

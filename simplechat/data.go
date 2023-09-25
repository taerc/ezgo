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
	connId string
}

type connectionContext struct {
	Id   string
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
	Id   string
	conn connection
}

type chatGroup struct {
	Id    string
	Admin string
	conn  connection

	userList *list.Element
}

type Client interface {
	NewClient(id string)
	GetId() string
	Login(id string) error
	SendMessageToUser(m Message) error
	SendMessageToGroup(m Message) error
}

type Group interface {
	NewGroup(id string)
	GetAdmin() error
	Login(id string) error
	GetId() string
	AddUserToGroup(usrId string) error
	RemoveUserFromGroup(usrId string) error
	GetUserList() *list.Element
	SendMessage(m Message) error
}

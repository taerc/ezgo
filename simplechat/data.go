package simplechat

import "container/list"

type Message struct {
	Id   string `json:"id"`
	From string `json:"from"`
	To   string `json:"to"`
	Type string `json:"type"`
	Data string `json:"data"`
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
	_user  map[string]*ChatUser
	_group map[string]*ChatGroup
)

func init() {
	_user = make(map[string]*ChatUser)
	_group = make(map[string]*ChatGroup)
}

type ChatUser struct {
	Id   string
	conn connection
}

type ChatGroup struct {
	Id    string
	Admin string
	conn  connection

	userList *list.Element
}

type Client interface {
	NewClient(id string) *ChatUser
	GetId() string
	Login(id string) error
	Logout(id string) error
	SendMessageToUser(m Message) error
	SendMessageToGroup(m Message) error
}

type Group interface {
	NewGroup(id string) *ChatGroup
	GetId() string
	GetAdmin() error
	Login(id string) error
	Logout(id string) error
	AddUserToGroup(usrId string) error
	RemoveUserFromGroup(usrId string) error
	GetUserList() *list.Element
	SendMessage(m Message) error
}

func trackUser(usr *ChatUser) {
	_user[usr.Id] = usr
}

func trackGroup(group *ChatGroup) {
	_group[group.Id] = group
}

package simplechat

import (
	"container/list"
	"sync"
)

type Command struct {
	Commd string      `json:"cmd"`
	Data  interface{} `json:"data"`
}

type Message struct {
	Type string `json:"type"`
	Id   string `json:"id"`
	From string `json:"from"`
	To   string `json:"to"`
	Data string `json:"data"`
}

type LoginMessage struct {
}

type LogoutMessage struct {
}

type JoinGroup struct {
}

type LeaveGroup struct {
}

type DestoryGroup struct {
}

type ChatUser struct {
	Id   string
	conn connection
}

type ChatGroup struct {
	Id    string
	Admin string
	conn  connection

	userList *list.List
	lockList *sync.Mutex
}

type Client interface {
	NewClient(id string) *ChatUser
	GetId() string
	SendMessageToUser(m Message) error
	SendMessageToGroup(m Message) error
}

type Group interface {
	NewGroup(id string) *ChatGroup
	GetId() string
	AddUserToGroup(usrId string) error
	RemoveUserFromGroup(usrId string) error
	GetUserList() *list.List
	SendMessage(m Message) error
}

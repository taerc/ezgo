package simplechat

import (
	"container/list"
	"sync"
)

type Message struct {
	Id   string `json:"id"`
	From string `json:"from"`
	To   string `json:"to"`
	Type string `json:"type"`
	Data string `json:"data"`
}

type connection struct {
	connId   string
	connLock sync.Mutex
}

type connectionContext struct {
	Id   string
	Type int
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

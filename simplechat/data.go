package simplechat

import (
	"container/list"
	"sync"
)

type Command struct {
	Cmd  string      `json:"cmd"`
	Data interface{} `json:"data"`
}

type Message struct {
	Type string `json:"type"`
	Id   string `json:"id"`
	From string `json:"from"`
	To   string `json:"to"`
	Data string `json:"data"`
}

type SendMessage struct {
	Type string `json:"type"`
	Id   string `json:"id"`
	From string `json:"from"`
	To   string `json:"to"`
	Data string `json:"data"`
}
type SendGroupMessage struct {
	Type    string `json:"type"`
	Id      string `json:"id"`
	From    string `json:"from"`
	GroupId string `json:"groupId"`
	Data    string `json:"data"`
}

type LoginMessage struct {
	UsrId string
}

type LogoutMessage struct {
	UsrId string
}

type JoinGroupMessage struct {
	UsrIds  []string
	GroupId string
}

type LeaveGroupMessage struct {
	UsrIds  []string
	GroupId string
}

type DestoryGroupMessage struct {
	GroupId string
}

type ChatUser struct {
	Id   string
	conn connection
}

type ChatGroup struct {
	Id    string
	Admin string

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

func NewCommand(cmd string, data interface{}) Command {

	return Command{
		Cmd:  cmd,
		Data: data,
	}

}

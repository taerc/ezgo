package simplechat

import (
	"container/list"
	"errors"
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

var (
	// _conn  map[string]connection
	_user      map[string]*ChatUser
	_lockUser  *sync.Mutex
	_group     map[string]*ChatGroup
	_lockGroup *sync.Mutex
)

func init() {
	_user = make(map[string]*ChatUser)
	_lockUser = &sync.Mutex{}
	_group = make(map[string]*ChatGroup)
	_lockGroup = &sync.Mutex{}
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

func trackUser(usr *ChatUser) {
	_user[usr.Id] = usr
}

func trackGroup(group *ChatGroup) {
	_group[group.Id] = group
}

func getUserById(usrId string) (*ChatUser, error) {
	_lockUser.Lock()
	defer _lockUser.Unlock()

	if c, ok := _user[usrId]; ok {
		return c, nil
	}
	return nil, errors.New("not found user")
}

func getGroupById(groupId string) (*ChatGroup, error) {
	_lockUser.Lock()
	defer _lockUser.Unlock()

	if g, ok := _group[groupId]; ok {
		return g, nil
	}
	return nil, errors.New("not found group")
}

package simplechat

import (
	"errors"
	"sync"
)

var (
	_connection     map[string]*connection
	_connectionLock *sync.Mutex
	_user           map[string]*ChatUser
	_lockUser       *sync.Mutex
	_group          map[string]*ChatGroup
	_lockGroup      *sync.Mutex
)

func init() {
	_connection = make(map[string]*connection)
	_connectionLock = &sync.Mutex{}
	_user = make(map[string]*ChatUser)
	_lockUser = &sync.Mutex{}
	_group = make(map[string]*ChatGroup)
	_lockGroup = &sync.Mutex{}
}

func trackUser(usr *ChatUser) {
	_user[usr.Id] = usr
}

func trackGroup(group *ChatGroup) {
	_group[group.Id] = group
}

func trackConnection(conn *connection) {
	_connection[conn.Id] = conn

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

func getConnectionById(connId string) (*connection, error) {
	_connectionLock.Lock()
	defer _connectionLock.Unlock()
	if c, ok := _connection[connId]; ok {
		return c, nil
	}
	return nil, errors.New("not found group")

}

package simplechat

import (
	"container/list"
	"fmt"
	"sync"
)

func NewGroup(id string) *ChatGroup {
	group := &ChatGroup{
		Id:       id,
		userList: list.New(),
		lockList: &sync.Mutex{},
		conn: connection{
			connLock: sync.Mutex{},
		},
	}
	trackGroup(group)
	return group
}

func (g *ChatGroup) GetId() string {
	return g.Id
}

func (g *ChatGroup) AddUserToGroup(usrId string) error {
	g.lockList.Lock()
	defer g.lockList.Unlock()
	g.userList.PushBack(usrId)
	return nil

}

func (g *ChatGroup) RemoveUserFromGroup(usrId string) error {
	g.lockList.Lock()
	defer g.lockList.Unlock()

	for item := g.userList.Front(); item != nil; item = item.Next() {
		u := item.Value.(string)
		if u == usrId {
			g.userList.Remove(item)
			break
		}
		fmt.Println(u)
	}
	return nil
}

func (g *ChatGroup) ShowGroup() {
	g.lockList.Lock()
	defer g.lockList.Unlock()
	for item := g.userList.Front(); item != nil; item = item.Next() {
		u := item.Value.(string)
		fmt.Println(u)
	}

}

func (g *ChatGroup) GetUserList() *list.List {
	return g.userList
}

func (g *ChatGroup) SendMessage(m Message) error {

	return nil
}

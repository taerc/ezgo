package simplechat

import (
	"encoding/json"
	"fmt"
	"sync"
)

func NewClient(id string) *ChatUser {

	usr := &ChatUser{
		Id: id,
		conn: connection{
			connLock: sync.Mutex{},
		},
	}
	trackUser(usr)
	return usr
}

func (c *ChatUser) GetId() string {
	return c.Id
}

func (c *ChatUser) SendMessageToUser(data string, usrId string) error {
	m := Message{
		From: c.Id,
		To:   usrId,
		Type: "single",
		Data: data,
	}

	msg, _ := json.MarshalIndent(m, " ", " ")
	c.conn.SendMessage(string(msg))
	return nil
}

func (c *ChatUser) SendMessageToGroup(data string, groudId string) error {

	group, e1 := getGroupById(groudId)

	if e1 != nil {
		fmt.Println(e1.Error())
		return e1
	}

	usrList := group.GetUserList()

	for item := usrList.Front(); item != nil; item = item.Next() {
		m := Message{
			From: c.Id,
			To:   groudId,
			Type: "group",
			Data: data,
		}
		usrId := item.Value.(string)
		if usrId == c.Id {
			continue
		}
		conn, _ := getUserById(usrId)
		msg, _ := json.MarshalIndent(m, " ", " ")
		conn.conn.SendMessage(string(msg))
	}

	return nil
}

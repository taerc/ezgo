package simplechat

import "encoding/json"

func NewClient(id string) *ChatUser {

	usr := &ChatUser{
		Id: id,
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

	m := Message{
		From: c.Id,
		To:   groudId,
		Type: "group",
		Data: data,
	}

	msg, _ := json.MarshalIndent(m, " ", " ")
	c.conn.SendMessage(string(msg))

	return nil
}

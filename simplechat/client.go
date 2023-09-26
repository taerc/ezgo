package simplechat

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

func (c *ChatUser) SendMessageToUser(m Message) error {

	return nil
}

func (c *ChatUser) SendMessageToGroup(m Message) error {

	return nil
}

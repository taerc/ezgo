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

func (c *ChatUser) Login(id string) error {
	return nil
}

func (c *ChatUser) Logout(id string) error {
	return nil
}

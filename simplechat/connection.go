package simplechat

import "fmt"

func (c *connection) SendMessage(data string) {
	c.connLock.Lock()
	defer c.connLock.Unlock()

	fmt.Println(c.connId, data)
}

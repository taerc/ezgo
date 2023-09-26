package simplechat

import "fmt"

func (c *connection) SendMessage(data string) {
	fmt.Println(c.connId, data)
}

package main

import (
	"github.com/taerc/ezgo/simplechat"
)

func GroupTesting() {
	// genId := ezgo.NewEZID(10, 10, ezgo.ChatIDSetting())
	// id, _ := genId.NextStringID()
	c1 := simplechat.NewClient("wangfangming")
	c2 := simplechat.NewClient("wang")
	c3 := simplechat.NewClient("fang")
	c4 := simplechat.NewClient("ming")

	g1 := simplechat.NewGroup("G1")

	g1.AddUserToGroup(c1.GetId())
	g1.AddUserToGroup(c2.GetId())
	g1.AddUserToGroup(c3.GetId())
	g1.AddUserToGroup(c4.GetId())

	// g1.RemoveUserFromGroup("ming")
	g1.ShowGroup()

	c1.SendMessageToUser("Hello ", "wang")

	c1.SendMessageToGroup("Hello", "G1")

}
func main() {
	// GroupTesting()
	simplechat.StartChatServer(9999)

}

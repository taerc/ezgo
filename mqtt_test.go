package ezgo

import (
	"fmt"
	"github.com/taerc/ezgo/conf"
	"testing"
	"time"
)

var c = conf.MQTTConf{
	MQTTSubTopic: "test.123456",
	MQTTHost:     "106.14.248.55",
	MQTTPort:     8899,
	MQTTUser:     "airia_pub",
	MQTTPwd:      "Airia12#$",
}

func process(message []byte) {
	fmt.Println(string(message))
}
func processAA(message []byte) {
	fmt.Println(string(message) + "AAAA")
}
func TestNewMQTTSubClient(t *testing.T) {

	mc := NewMQTTSubClient("sub", &c)
	mc.AddConsumer(c.MQTTSubTopic, process)
	mc.AddConsumer("AAAA", processAA)

	for {
		Info(nil, "TESTING", "sub")
		time.Sleep(time.Minute * 5)
	}
}

func TestNewMQTTPubClient(t *testing.T) {
	mc := NewMQTTPubClient("pub", &c)
	mc.Publish(c.MQTTSubTopic, "{123456}")
	mc.Publish("CCCC", "{123456}")
	mc.Publish("AAAA", "{8123456}")

}

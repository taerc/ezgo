package ezgo

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/taerc/ezgo/conf"
	"time"
)

type MQTTClient struct {
	client   mqtt.Client
	Name     string
	Consumer map[string]func(message []byte)
}

func newMQTTClient() *MQTTClient {
	mc := new(MQTTClient)
	mc.Consumer = make(map[string]func(message []byte))
	return mc
}

func NewMQTTSubClient(name string, conf *conf.MQTTConf) *MQTTClient {

	mc := newMQTTClient()
	mc.Name = name
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", conf.MQTTHost, conf.MQTTPort))
	opts.SetClientID(name)
	opts.SetUsername(conf.MQTTUser)
	opts.SetPassword(conf.MQTTPwd)

	opts.SetDefaultPublishHandler(mc.defaultSubClient)
	opts.OnConnect = mc.defaultOnConnect
	opts.OnConnectionLost = mc.defaultOnConnectionLost

	mc.client = mqtt.NewClient(opts)
	if token := mc.client.Connect(); token.Wait() && token.Error() != nil {
		time.Sleep(time.Second)
		return nil
	}

	return mc
}

func (mc *MQTTClient) AddConsumer(topic string, cons func(message []byte)) {
	mc.Consumer[topic] = cons
	mc.client.Subscribe(topic, 1, nil)
}

func (mc *MQTTClient) defaultSubClient(client mqtt.Client, msg mqtt.Message) {

	topic := msg.Topic()
	Info(nil, M, topic)
	if c, ok := mc.Consumer[topic]; ok {
		c(msg.Payload())
	}
}

func (mc *MQTTClient) defaultOnConnect(client mqtt.Client) {
	Info(nil, M, "connecting")
}

func (mc *MQTTClient) defaultOnConnectionLost(client mqtt.Client, e error) {
	Info(nil, M, e.Error())
}

func NewMQTTPubClient(name string, conf *conf.MQTTConf) *MQTTClient {
	mc := newMQTTClient()
	mc.Name = name
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", conf.MQTTHost, conf.MQTTPort))
	opts.SetClientID(name)
	opts.SetUsername(conf.MQTTUser)
	opts.SetPassword(conf.MQTTPwd)

	opts.OnConnect = mc.defaultOnConnect
	opts.OnConnectionLost = mc.defaultOnConnectionLost

	mc.client = mqtt.NewClient(opts)
	if token := mc.client.Connect(); token.Wait() && token.Error() != nil {
		time.Sleep(time.Second)
		return nil
	}

	return mc
}

func (mc *MQTTClient) Publish(topic string, data any) {
	mc.client.Publish(topic, 1, false, data)
}

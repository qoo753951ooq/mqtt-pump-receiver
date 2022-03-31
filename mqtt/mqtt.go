package mqtt

import (
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type MqttConfig struct {
	Ip       string
	Port     int
	User     string
	Pwd      string
	Topic    string
	Qos      byte
	ClientID string
}

type Client struct {
	client mqtt.Client
}

var Choke = make(chan [2]string)

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	Choke <- [2]string{msg.Topic(), string(msg.Payload())}
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("Connected")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("Connect lost: %v", err)
}

func NewClientOptions(conf MqttConfig) *mqtt.ClientOptions {

	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tls://%s:%d", conf.Ip, conf.Port))
	opts.SetUsername(conf.User)
	opts.SetPassword(conf.Pwd)

	opts.ClientID = conf.ClientID
	opts.DefaultPublishHandler = messagePubHandler
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler

	return opts
}

func NewClient(opts *mqtt.ClientOptions) (Client, error) {

	c := Client{}

	c.client = mqtt.NewClient(opts)

	if token := c.client.Connect(); token.Wait() && token.Error() != nil {
		fmt.Printf("token err: %s", token.Error())
		return c, token.Error()
	}

	return c, nil
}

func (c Client) Subscribe(conf MqttConfig) error {

	if token := c.client.Subscribe(conf.Topic, conf.Qos, nil); token.Wait() && token.Error() != nil {
		fmt.Printf("subscribe err: %s", token.Error())
		return token.Error()
	}

	return nil
}

func (c Client) Disconnect(quiesce uint) {
	c.client.Disconnect(quiesce)
}

func (c Client) IsConnected() bool {
	return c.client.IsConnected()
}

package main

import (
	"fmt"
	"mqtt-pump-receiver/conf"
	"mqtt-pump-receiver/mqtt"
	"mqtt-pump-receiver/service"
	"mqtt-pump-receiver/util"

	"github.com/spf13/viper"
)

func main() {

	conf.Init()
	writeLog := viper.GetBool("log")
	mqttConf := createMqttConf()

	opts := mqtt.NewClientOptions(mqttConf)

	pumpStMap := service.GetPumpStMap(viper.GetString("pump.org_id"),
		viper.GetString("pump.display"), viper.GetString("pump.etype"))

creatMqttCli:

	client, err := mqtt.NewClient(opts)

	if err != nil {
		fmt.Printf("%s:%s \n", "mqtt New Client", err.Error())
		return
	}

	if err := client.Subscribe(mqttConf); err != nil {
		fmt.Printf("%s:%s \n", "mqtt Client Subscribe", err.Error())
		return
	}

	for {

		incoming := <-mqtt.Choke

		if writeLog {
			util.Writelog(fmt.Sprintf("MESSAGE: %s", incoming[1]))
		}

		sendDatas := service.GetSendPumpDatas(incoming[1], pumpStMap)

		if len(sendDatas) > 0 {
			service.PostPumpDataToCOVM(sendDatas)
		}

		if !client.IsConnected() {
			fmt.Println("client Disconnected")
			client.Disconnect(60)
			goto creatMqttCli
		}
	}
}

func createMqttConf() mqtt.MqttConfig {

	return mqtt.MqttConfig{
		Ip:       viper.GetString("mqtt.ip"),
		Port:     viper.GetInt("mqtt.port"),
		User:     viper.GetString("mqtt.user"),
		Pwd:      viper.GetString("mqtt.pwd"),
		Topic:    viper.GetString("mqtt.topic"),
		Qos:      1,
		ClientID: util.NewV4UUID(),
	}
}

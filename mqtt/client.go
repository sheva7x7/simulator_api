package mqtt

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var client mqtt.Client

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("Connected")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("Connect lost: %v", err)
}

func StartClient() {
	fmt.Println("starting mqtt client")
	broker := os.Getenv("MQTT_URL")
	port, _ := strconv.Atoi(os.Getenv("MQTT_PORT"))
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", broker, port))
	opts.SetDefaultPublishHandler(messagePubHandler)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler
	client = mqtt.NewClient(opts)
	token := client.Connect()
	if token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	tokenobj, _ := json.Marshal(token)
	fmt.Println(string(tokenobj))
}

func Publish(topic string, payload interface{}) {
	payloadString, _ := json.Marshal(payload)
	fmt.Println("pubbing", string(payloadString))
	token := client.Publish(topic, 0, false, string(payloadString))
	token.Wait()
}

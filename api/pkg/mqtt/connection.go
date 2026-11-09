package mqtt

import (
	"fmt"
	"log"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func NewMqttClient(brokenUrl string) (mqtt.Client, error) {
	options := mqtt.NewClientOptions()
	options.AddBroker(brokenUrl)
	options.SetClientID("water_iot_api_subcriber")
	options.SetAutoReconnect(true)
	options.SetConnectRetry(true)
	options.SetConnectRetryInterval(2 * time.Second)

	options.OnConnect = func(c mqtt.Client) {
		log.Println("[MQTT] Connected to Broken Successfully!")
	}

	options.OnConnectionLost = func(c mqtt.Client, err error) {
		log.Printf("[MQTT] Connection lost: %v", err)
	}

	client := mqtt.NewClient(options)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return nil, fmt.Errorf("failed to connect to mqtt broker: %w", token.Error())
	}

	return client, nil
}

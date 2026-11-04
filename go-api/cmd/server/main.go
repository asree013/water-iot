package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"water-iot-api/cmd/models"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var (
	latestData = make(map[string]models.IoTData)
	dataMutex  sync.RWMutex
)

func onMessageReceived(client mqtt.Client, message mqtt.Message) {
	var payload models.IoTData
	err := json.Unmarshal(message.Payload(), &payload)
	if err != nil {
		log.Printf("Failed to parse JSON: %v\n", err)
		return
	}

	key := fmt.Sprintf("%s/%s", payload.Zone, payload.Pipe)
	dataMutex.Lock()
	latestData[key] = payload
	dataMutex.Unlock()

	log.Printf("[API MQTT Receiver] Recv from %s | Pressure: %.2f | Flow: %.2f",
		message.Topic(), payload.Pressure, payload.FlowRate)
}

func connectMQTT() mqtt.Client {
	brokerURL := os.Getenv("MQTT_BROKER")
	if brokerURL == "" {
		brokerURL = "tcp://localhost:1883"
	}

	opts := mqtt.NewClientOptions()
	opts.AddBroker(brokerURL)
	opts.SetClientID("water_iot_go_api_service")
	opts.SetAutoReconnect(true)
	opts.SetConnectRetry(true)

	opts.OnConnect = func(c mqtt.Client) {
		log.Println("MQTT Connected! Subscribing to topic: water-iot/sensors/#")
		go func() {
			topic := "water-iot/sensors/#"
			if token := c.Subscribe(topic, 0, onMessageReceived); token.Wait() && token.Error() != nil {
				log.Printf("Warning: Failed to subscribe: %v", token.Error())
			} else {
				log.Println("Successfully subscribed to water-iot/sensors/#")
			}
		}()
	}

	opts.OnConnectionLost = func(c mqtt.Client, err error) {
		log.Printf("MQTT Connection lost: %v", err)
	}

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("MQTT Connect error: %v", token.Error())
	}
	return client
}

func main() {
	mqttClient := connectMQTT()
	defer mqttClient.Disconnect(250)

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Go IoT Ingestion API is running!")
	})

	http.HandleFunc("/api/sensors/latest", func(w http.ResponseWriter, r *http.Request) {
		dataMutex.RLock()
		defer dataMutex.RUnlock()

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(latestData)
	})

	fmt.Println("HTTP Server started at http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("HTTP server error: %v", err)
	}
}

package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type IoTData struct {
	Pressure         float64 `json:"pressure"`
	FlowRate         float64 `json:"flow_rate"`
	Temperature      float64 `json:"temperature"`
	Vibration        float64 `json:"vibration"`
	RPM              float64 `json:"rpm"`
	OperationalHours int     `json:"operational_hours"`
	Zone             string  `json:"zone"`
	Block            string  `json:"block"`
	Pipe             string  `json:"pipe"`
	LocationCode     string  `json:"location_code"`
	Latitude         float64 `json:"latitude"`
	Longitude        float64 `json:"longitude"`
	LeakageFlag      int     `json:"leakage_flag"`
	Timestamp        int64   `json:"timestamp"` // Unix timestamp เพิ่มความสมจริงของ IoT
}

func parseFloat(s string) float64 {
	val, _ := strconv.ParseFloat(s, 64)
	return val
}

func parseInt(s string) int {
	val, _ := strconv.Atoi(s)
	return val
}

func main() {
	filePath := "data/location_aware_gis_leakage_dataset.csv"
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Failed to open CSV file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatalf("Failed to read CSV records: %v", err)
	}

	dataRows := records[1:]
	fmt.Printf("Loaded %d rows from dataset. Starting MQTT IoT simulation...\n", len(dataRows))

	// ดึง Broker URL จาก Environment Variable ถ้าไม่มีให้ใช้ localhost:1883
	brokerURL := os.Getenv("MQTT_BROKER")
	if brokerURL == "" {
		brokerURL = "tcp://localhost:1883"
	}

	opts := mqtt.NewClientOptions()
	opts.AddBroker(brokerURL)
	opts.SetClientID("water_iot_stream_publisher")
	opts.SetAutoReconnect(true)
	opts.SetConnectRetry(true)
	opts.SetConnectRetryInterval(2 * time.Second)

	client := mqtt.NewClient(opts)
	fmt.Printf("Connecting to MQTT Broker at %s ...\n", brokerURL)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("Failed to connect to MQTT broker: %v", token.Error())
	}
	defer client.Disconnect(250)
	fmt.Println("Connected to MQTT Broker successfully!")

	ticker := time.NewTicker(1 * time.Second) // ส่งข้อมูลทุกๆ 1 วินาที
	defer ticker.Stop()

	rowIndex := 0
	for range ticker.C {
		if rowIndex >= len(dataRows) {
			fmt.Println("Reached end of dataset. Looping back to start...")
			rowIndex = 0
		}

		row := dataRows[rowIndex]
		rowIndex++

		data := IoTData{
			Pressure:         parseFloat(row[0]),
			FlowRate:         parseFloat(row[1]),
			Temperature:      parseFloat(row[2]),
			Vibration:        parseFloat(row[3]),
			RPM:              parseFloat(row[4]),
			OperationalHours: parseInt(row[5]),
			Zone:             row[6],
			Block:            row[7],
			Pipe:             row[8],
			LocationCode:     row[9],
			Latitude:         parseFloat(row[10]),
			Longitude:        parseFloat(row[11]),
			LeakageFlag:      parseInt(row[12]),
			Timestamp:        time.Now().Unix(),
		}

		jsonPayload, err := json.Marshal(data)
		if err != nil {
			log.Printf("JSON marshal error: %v", err)
			continue
		}

		// Topic ตามมาตรฐาน IoT: แบ่งตาม Zone และ Pipe
		// เช่น water-iot/sensors/Zone_A/Pipe_01
		topic := fmt.Sprintf("water-iot/sensors/%s/%s", data.Zone, data.Pipe)

		// QoS 0 หรือ 1 (ใน IoT นิยม QoS 0 สำหรับ telemetry streaming)
		token := client.Publish(topic, 0, false, jsonPayload)
		token.Wait()

		fmt.Printf("[MQTT Sent] Row %d -> Topic: %s | Pressure: %.2f | Flow: %.2f\n",
			rowIndex, topic, data.Pressure, data.FlowRate)
	}
}

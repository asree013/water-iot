package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"time"
)

// หมายเหตุ tag json: ไม่ควรมีเว้นวรรคหลังเครื่องหมายโคลอน เช่น `json:"pressure"`
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
	fmt.Printf("Loaded %d rows from dataset. Starting UDP stream...\n", len(dataRows))

	// 1. กำหนด IP และ Port ปลายทางของ UDP (เช่น localhost:7672 หรือ service_name ใน Docker)
	targetAddr := "127.0.0.1:7672"
	udpAddr, err := net.ResolveUDPAddr("udp", targetAddr)
	if err != nil {
		log.Fatalf("Failed to resolve UDP address: %v", err)
	}

	// 2. สร้าง UDP Connection
	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		log.Fatalf("Failed to dial UDP: %v", err)
	}
	defer conn.Close()

	fmt.Printf("Streaming UDP packets to %s ...\n", targetAddr)

	ticker := time.NewTicker(500 * time.Millisecond) // ปรับความเร็วในการพ่นข้อมูลได้ตามต้องการ เช่น ทุกๆ 500ms
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
		}

		jsonValue, err := json.Marshal(data)
		if err != nil {
			log.Printf("JSON marshal error: %v", err)
			continue
		}

		// 3. พ่นแพ็กเก็ต UDP ออกไป (Fire-and-forget)
		_, err = conn.Write(jsonValue)
		if err != nil {
			log.Printf("Error sending UDP packet: %v", err)
			continue
		}

		fmt.Printf("[UDP Stream] Sent row %d -> Zone: %s, Pipe: %s (%d bytes)\n",
			rowIndex, data.Zone, data.Pipe, len(jsonValue))
	}
}

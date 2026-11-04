package models

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

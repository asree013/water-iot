package main

import (
	"api/config"
	"api/modules/entities"
	"api/modules/servers"
	"api/pkg/database"
	"log"
)

// @title           Water IoT API
// @version         1.0
// @description     API documentation for Water IoT System
// @host            localhost:8080
// @BasePath        /
func main() {
	// load .env
	cfg := config.LoadConfig()

	// connect mysql database
	db, err := database.NewSQLConnection(cfg)
	if err != nil {
		log.Fatalf("Database connection failded: %v", err)
	}
	// auto migrate table
	if err := db.AutoMigrate(&entities.WaterData{}, &entities.TDMAZone{}); err != nil {
		log.Printf("AutoMigrate warning: %v", err)
	}

	server := servers.NewServers(cfg, db)
	server.Start()
}

// func main() {
// 	// 1. Load env (รองรับทั้งรันจากโฟลเดอร์ api/ หรือ api/apps/)
// 	_ = godotenv.Load(".env", "../.env")

// 	cfg := new(config.Configs)

// 	// Fiber configs
// 	cfg.App.Host = getEnv("APP_HOST", "0.0.0.0")
// 	cfg.App.Port = getEnv("APP_PORT", "8080")
// 	cfg.App.Env = getEnv("APP_ENV", "development")

// 	// Database configs
// 	maxOpenConns, _ := strconv.Atoi(getEnv("DB_MAX_OPEN_CONNS", "25"))
// 	maxIdleConns, _ := strconv.Atoi(getEnv("DB_MAX_IDLE_CONNS", "10"))

// 	cfg.MySQL = config.MySQLConfig{
// 		Host:            getEnv("DB_HOST", "localhost"),
// 		Port:            getEnv("DB_PORT", "3306"),
// 		User:            getEnv("DB_USER", "root"),
// 		Password:        getEnv("DB_PASSWORD", "12345678"),
// 		DBName:          getEnv("DB_NAME", "water_iot_db"),
// 		MaxOpenConns:    maxOpenConns,
// 		MaxIdleConns:    maxIdleConns,
// 		ConnMaxLifetime: time.Hour,
// 	}

// 	// MQTT configs
// 	cfg.MQTT.BrokerURL = getEnv("MQTT_BROKER", "tcp://localhost:1883")

// 	// 2. เชื่อมต่อ MySQL Database
// 	db, err := database.NewMySQLConnection(cfg)
// 	if err != nil {
// 		log.Fatalf("Database connection failed: %v", err)
// 	}

// 	// 3. Auto Migrate (สร้าง Table ให้อัตโนมัติถ้ายังไม่มี)
// 	_ = db.AutoMigrate(&entities.WaterData{})

// 	// 4. สตาร์ท Server
// 	server := servers.NewServer(cfg, db)
// 	server.Start()
// }

// func getEnv(key, defaultVal string) string {
// 	if val, ok := os.LookupEnv(key); ok && val != "" {
// 		return val
// 	}
// 	return defaultVal
// }

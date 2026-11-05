package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Configs struct {
	App   AppConfig
	MySQL MySQLConfig
	MQTT  MQTTConfig
}

type AppConfig struct {
	Port string
	Env  string
}

type MySQLConfig struct {
	Host            string
	Port            string
	User            string
	Password        string
	DBName          string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

type MQTTConfig struct {
	BrokerURL string
}

// DSN สร้าง Data Source Name สำหรับ MySQL
func (m *MySQLConfig) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		m.User,
		m.Password,
		m.Host,
		m.Port,
		m.DBName,
	)
}

// ConnectGorm ทำการเชื่อมต่อ MySQL ด้วย GORM ORM พร้อมตั้งค่า Connection Pool
func (m *MySQLConfig) ConnectGorm() (*gorm.DB, error) {
	dsn := m.DSN()

	logMode := logger.Info
	if os.Getenv("APP_ENV") == "production" {
		logMode = logger.Warn
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logMode),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MySQL: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB from gorm: %w", err)
	}

	// กำหนด Connection Pool
	sqlDB.SetMaxOpenConns(m.MaxOpenConns)
	sqlDB.SetMaxIdleConns(m.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(m.ConnMaxLifetime)

	return db, nil
}

// LoadConfig โหลด configuration จาก .env หรือ System Environment
func LoadConfig() *Configs {
	if err := godotenv.Load(); err != nil {
		log.Println("Notice: .env file not found or using system environment variables")
	}

	return &Configs{
		App: AppConfig{
			Port: getEnv("APP_PORT", "8080"),
			Env:  getEnv("APP_ENV", "development"),
		},
		MySQL: MySQLConfig{
			Host:            getEnv("DB_HOST", "localhost"),
			Port:            getEnv("DB_PORT", "3306"),
			User:            getEnv("DB_USER", "root"),
			Password:        getEnv("DB_PASSWORD", "rootpassword"),
			DBName:          getEnv("DB_NAME", "water_iot_db"),
			MaxOpenConns:    getEnvAsInt("DB_MAX_OPEN_CONNS", 25),
			MaxIdleConns:    getEnvAsInt("DB_MAX_IDLE_CONNS", 10),
			ConnMaxLifetime: time.Hour,
		},
		MQTT: MQTTConfig{
			BrokerURL: getEnv("MQTT_BROKER", "tcp://localhost:1883"),
		},
	}
}

func getEnv(key, defaultVal string) string {
	if val, exists := os.LookupEnv(key); exists && val != "" {
		return val
	}
	return defaultVal
}

func getEnvAsInt(key string, defaultVal int) int {
	valStr := getEnv(key, "")
	if val, err := strconv.Atoi(valStr); err == nil {
		return val
	}
	return defaultVal
}

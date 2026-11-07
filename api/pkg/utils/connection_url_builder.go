package utils

import (
	"api/config"
	"errors"
	"fmt"
)

func ConnectionUrlBuilder(stuff string, cfg *config.Configs) (string, error) {
	var url string

	switch stuff {
	case "fiber":
		url = fmt.Sprintf("%s:%s", cfg.App.Host, cfg.App.Port)
	case "mysql":
		url = fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			cfg.MySQL.User,
			cfg.MySQL.Password,
			cfg.MySQL.Host,
			cfg.MySQL.Port,
			cfg.MySQL.DBName,
		)
	case "mqtt":
		url = cfg.MQTT.BrokerURL
	default:
		errMsg := fmt.Sprintf("error, connection url builder doesn't know the %s", stuff)
		return "", errors.New(errMsg)
	}

	return url, nil
}

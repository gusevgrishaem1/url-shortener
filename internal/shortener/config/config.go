package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

// Config Конфигурация приложения
type Config struct {
	Port        string // Port порт, где будет запущенно приложение.
	ServerURL   string // ServerURL адрес web сервиса, который сократил ссылку.
	Database    string // Database представляет собой URL базы данных, используемой приложением.
	EnableHTTPS bool   // EnableHTTPS представляет собой флаг, указывающий на включение HTTPS сервера.
}

// NewConfig инициализирует конфигурацию из переменных среды.
func NewConfig() *Config {
	for _, env := range os.Environ() {
		fmt.Println(env)
	}
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080" // Значение по умолчанию
	}

	serverURL := os.Getenv("SERVER_URL")
	if serverURL == "" {
		log.Fatal("ServerURL is required")
	}

	database := os.Getenv("DATABASE_URL")
	if database == "" {
		log.Fatal("DATABASE_URL is required")
	}

	enableHTTPSStr := os.Getenv("ENABLE_HTTPS")
	enableHTTPS, err := strconv.ParseBool(enableHTTPSStr)
	if err != nil {
		enableHTTPS = false // Значение по умолчанию
	}

	return &Config{Port: port, Database: database, EnableHTTPS: enableHTTPS, ServerURL: serverURL}
}

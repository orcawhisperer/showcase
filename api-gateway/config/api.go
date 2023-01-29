package config

import (
	"os"
)

type server struct {
	HTTPHost string
	HTTPPort string
}

type logger struct {
	Level string
}

// Settings struct
type Settings struct {
	Server *server
	Logger *logger
}

// GetSettings returns the settings
func GetSettings() *Settings {

	Settings := &Settings{
		Server: &server{
			HTTPHost: os.Getenv("HTTP_HOST"),
			HTTPPort: os.Getenv("HTTP_PORT"),
		},

		Logger: &logger{
			Level: os.Getenv("LOG_LEVEL"),
		},
	}
	return Settings
}

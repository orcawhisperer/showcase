package config

import (
	"os"
)

type server struct {
	GrpcHost string
	GrcpPort string
	HTTPHost string
	HTTPPort string
}

type database struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SslMode  string
}

type logger struct {
	Level string
}

// Settings struct
type Settings struct {
	Server   *server
	Database *database
	Logger   *logger
}

// GetSettings returns the settings
func GetSettings() *Settings {

	Settings := &Settings{
		Server: &server{
			GrpcHost: os.Getenv("VIDEO_SVC_GRPC_HOST"),
			GrcpPort: os.Getenv("VIDEO_SVC_GRPC_PORT"),
			HTTPHost: os.Getenv("VIDEO_SVC_HTTP_HOST"),
			HTTPPort: os.Getenv("VIDEO_SVC_HTTP_PORT"),
		},

		Database: &database{
			Host:     os.Getenv("VIDEO_SVC_DB_HOST"),
			Port:     os.Getenv("VIDEO_SVC_DB_PORT"),
			User:     os.Getenv("VIDEO_SVC_DB_USER"),
			Password: os.Getenv("VIDEO_SVC_DB_PASSWORD"),
			Name:     os.Getenv("VIDEO_SVC_DB_NAME"),
			SslMode:  os.Getenv("VIDEO_SVC_DB_SSLMODE"),
		},

		Logger: &logger{
			Level: os.Getenv("VIDEO_SVC_LOG_LEVEL"),
		},
	}
	return Settings
}

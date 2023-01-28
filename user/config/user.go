package config

import (
	"os"
	"strconv"
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

type jwt struct {
	Secret string
	Expiry int
}

// Settings struct
type Settings struct {
	Server   *server
	Database *database
	Logger   *logger
	JWT      *jwt
}

// GetSettings returns the settings
func GetSettings() *Settings {

	jwt_expiry, _ := strconv.Atoi(os.Getenv("JWT_EXPIRY"))

	Settings := &Settings{
		Server: &server{
			GrpcHost: os.Getenv("GRPC_HOST"),
			GrcpPort: os.Getenv("GRPC_PORT"),
			HTTPHost: os.Getenv("HTTP_HOST"),
			HTTPPort: os.Getenv("HTTP_PORT"),
		},

		Database: &database{
			Host:     os.Getenv("DB_HOST"),
			Port:     os.Getenv("DB_PORT"),
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			Name:     os.Getenv("DB_NAME"),
			SslMode:  os.Getenv("DB_SSLMODE"),
		},

		Logger: &logger{
			Level: os.Getenv("LOG_LEVEL"),
		},
		JWT: &jwt{
			Secret: os.Getenv("JWT_SECRET"),
			Expiry: jwt_expiry,
		},
	}
	return Settings
}

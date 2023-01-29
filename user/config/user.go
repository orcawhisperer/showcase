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

	jwt_expiry, _ := strconv.Atoi(os.Getenv("USER_SVC_JWT_EXPIRY"))

	Settings := &Settings{
		Server: &server{
			GrpcHost: os.Getenv("USER_SVC_GRPC_HOST"),
			GrcpPort: os.Getenv("USER_SVC_GRPC_PORT"),
			HTTPHost: os.Getenv("USER_SVC_HTTP_HOST"),
			HTTPPort: os.Getenv("USER_SVC_HTTP_PORT"),
		},

		Database: &database{
			Host:     os.Getenv("USER_SVC_DB_HOST"),
			Port:     os.Getenv("USER_SVC_DB_PORT"),
			User:     os.Getenv("USER_SVC_DB_USER"),
			Password: os.Getenv("USER_SVC_DB_PASSWORD"),
			Name:     os.Getenv("USER_SVC_DB_NAME"),
			SslMode:  os.Getenv("USER_SVC_DB_SSLMODE"),
		},

		Logger: &logger{
			Level: os.Getenv("USER_SVC_LOG_LEVEL"),
		},
		JWT: &jwt{
			Secret: os.Getenv("USER_SVC_JWT_SECRET"),
			Expiry: jwt_expiry,
		},
	}
	return Settings
}

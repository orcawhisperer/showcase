package utils

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func GetDBConnection(dsn string) (*gorm.DB, error) {
	//return gorm connection using connectoin string
	conn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return conn, nil
}

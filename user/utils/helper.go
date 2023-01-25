package utils

import (
	"github.com/jinzhu/gorm"
)

func GetDBConnection(conn_str string) (*gorm.DB, error) {
	//return gorm connection using connectoin string
	conn, err := gorm.Open("postgres", conn_str)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

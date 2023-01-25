// User data model

package model

import (
	"regexp"

	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	EmailRegex = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)

	PhoneRegex = regexp.MustCompile(`^[0-9]{10}$`)

	PasswordRegex = regexp.MustCompile(`^[a-zA-Z0-9!@#$%^&*]{8,30}$`)

	NameRegex = regexp.MustCompile(`^[a-zA-Z ]+$`)

	IDRegex = regexp.MustCompile(`^[a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12}$`)
)

type User struct {
	Id       string `gorm:"primary_key"`
	Name     string
	Email    string `gorm:"unique_index"`
	Phone    string `gorm:"unique_index"`
	Password string
}

// Hook before create to generate uuid
func (u *User) BeforeCreate(tx *gorm.DB) error {
	uuid := uuid.NewV4()
	u.Id = uuid.String()
	return nil
}

// Hook before save to hash password using bcrypt

func (u *User) BeforeSave() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

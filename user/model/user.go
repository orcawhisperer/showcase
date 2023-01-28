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

	PasswordRegex = regexp.MustCompile(`^[a-zA-Z0-9!@#$%^&*]{8,}$`)

	NameRegex = regexp.MustCompile(`^[a-zA-Z ]+$`)

	IDRegex = regexp.MustCompile(`^[a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12}$`)
)

type User struct {
	gorm.Model
	UUID      string `gorm:"primary_key" json:"id"`
	FirstName string `gorm:"not null" json:"first_name"`
	LastName  string `gorm:"not null" json:"last_name"`
	Username  string `gorm:"unique_index" json:"username"`
	Email     string `gorm:"unique_index" json:"email"`
	Phone     string `gorm:"unique_index" json:"phone"`
	Password  string `gorm:"not null" json:"-"`
}

// Hook before create to generate uuid and hash password
func (u *User) BeforeCreate(tx *gorm.DB) error {
	uuid := uuid.NewV4()
	u.UUID = uuid.String()
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

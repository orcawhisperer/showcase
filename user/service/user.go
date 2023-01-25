// user service package

package service

import (
	"log"
	"os"
)

type User struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

type UserService struct {
	repo *UserRepo
	log  *log.Logger
}

func NewUserService(repo *UserRepo) *UserService {
	return &UserService{
		repo: repo,
		log:  log.New(os.Stdout, "user-service: ", log.LstdFlags),
	}
}

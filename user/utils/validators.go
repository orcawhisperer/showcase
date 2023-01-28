package utils

import (
	"fmt"

	pb "github.com/iamvasanth07/showcase/common/protos/user"
	"github.com/iamvasanth07/showcase/user/model"
)

// ValidateEmail validates email
func ValidateEmail(email string) error {
	if email == "" {
		return fmt.Errorf("email is required")
	}
	if !model.EmailRegex.MatchString(email) {
		return fmt.Errorf("invalid email")
	}
	if len(email) > 100 {
		return fmt.Errorf("email is too long")
	}
	return nil
}

// ValidatePassword validates password
func ValidatePassword(password string) error {
	if password == "" {
		return fmt.Errorf("password is required")
	}
	if len(password) < 8 {
		return fmt.Errorf("password is too short")
	}
	if !model.PasswordRegex.MatchString(password) {
		return fmt.Errorf("password is weak")
	}
	return nil
}

// ValidateUsername validates username
func ValidateUserName(username string) error {
	if username == "" {
		return fmt.Errorf("username is required")
	}
	if len(username) < 3 {
		return fmt.Errorf("username is too short")
	}
	if len(username) > 100 {
		return fmt.Errorf("username is too long")
	}
	return nil
}

// ValidateName validates Firstname and Lastname
func ValidateName(firstName string, lastName string) error {
	if firstName == "" {
		return fmt.Errorf("firstname is required")
	}
	if lastName == "" {
		return fmt.Errorf("lastname is required")
	}
	return nil
}

// ValidatePhone validates phone
func ValidatePhone(phone string) error {
	if phone == "" {
		return fmt.Errorf("phone is required")
	}
	if !model.PhoneRegex.MatchString(phone) {
		return fmt.Errorf("invalid phone")
	}
	return nil
}

// ValidateID validates id
func ValidateID(id string) error {
	if id == "" {
		return fmt.Errorf("id is required")
	}
	if !model.IDRegex.MatchString(id) {
		return fmt.Errorf("invalid id")
	}
	return nil
}

// ValidateUserCreate validates user
func ValidateUserCreate(user *pb.User) error {
	if user == nil {
		return fmt.Errorf("user is required")
	}
	if err := ValidateEmail(user.Email); err != nil {
		return err
	}
	if err := ValidatePassword(user.Password); err != nil {
		return err
	}
	if err := ValidateUserName(user.Username); err != nil {
		return err
	}
	if err := ValidateName(user.FirstName, user.LastName); err != nil {
		return err
	}
	if err := ValidatePhone(user.Phone); err != nil {
		return err
	}
	return nil
}

// ValidateUserUpdate validates user update
func ValidateUserUpdate(user *pb.User) error {

	if user == nil {
		return fmt.Errorf("user is required")
	}
	if err := ValidateEmail(user.Email); err != nil {
		return err
	}
	if err := ValidateName(user.FirstName, user.LastName); err != nil {
		return err
	}
	if err := ValidatePhone(user.Phone); err != nil {
		return err
	}
	return nil
}

// ValidateUserDelete validates user delete
func ValidateUserDelete(id string) error {
	if err := ValidateID(id); err != nil {
		return err
	}
	return nil
}

// ValidateUserGet validates user get
func ValidateUserGet(id string) error {
	if err := ValidateID(id); err != nil {
		return err
	}
	return nil
}

// ValidateUserGetAll validates user get all
func ValidateUserGetAll(req *pb.GetAllUserRequest) error {
	if req == nil {
		return fmt.Errorf("body cannot be empty")
	}
	if req.Paginate == nil {
		return fmt.Errorf("paginate options cannot be empty")
	}
	if req.Paginate.Limit < 0 {
		return fmt.Errorf("invalid limit")
	}
	if req.Paginate.Page < 0 {
		return fmt.Errorf("invalid page")
	}
	return nil
}

// ValidateUserLogin validates user login
func ValidateUserLogin(req *pb.LoginRequest) error {
	if req == nil {
		return fmt.Errorf("body cannot be empty")
	}
	if err := ValidateEmail(req.Email); err != nil {
		return err
	}
	if err := ValidatePassword(req.Password); err != nil {
		return err
	}
	return nil
}

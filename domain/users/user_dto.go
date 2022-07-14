package users

import (
	"bookstore_api/utils/errors"
	"strings"
)

const (
	StatusActive = "active"
)

type User struct {
	Id          int64  `json:"id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	DateCreated string `json:"date_created"`
	Status      string `json:"status"`
	Password    string `json:"password"`
}

func (user *User) Validate() *errors.RestErr {
	user.FirstName = strings.TrimSpace(user.FirstName)
	user.LastName = strings.TrimSpace(user.LastName)
	user.Password = strings.TrimSpace(user.Password)
	user.Email = strings.TrimSpace(strings.ToLower(user.Email))
	if user.Email == "" {
		return errors.NewBadRequestError("Invalid Email Address")
	}
	if user.Password == "" {
		return errors.NewBadRequestError("Invalid Password")
	}
	return nil
}

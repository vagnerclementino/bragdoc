package domain

import (
	"errors"
	"fmt"
	"net/mail"

	"github.com/google/uuid"
)

type User struct {
	ID    string
	Name  string
	Email string
}

func (u *User) Validate() error {

	if _, err := uuid.Parse(u.ID); err != nil {
		return errors.New("User.ID: the user's id is not valid")
	}

	if u.Name == "" {
		return errors.New("User.Name: the user's name cannot be empty")
	}

	if len(u.Name) <= 3 {
		return errors.New("User.Name: the user's name has an unexpected size")
	}

	if u.Email == "" {
		return errors.New("User.Email: the user's email cannot be empty")
	}

	if _, err := mail.ParseAddress(u.Email); err != nil {
		return fmt.Errorf("User.Email: the user's email %s is not valid", u.Email)
	}

	return nil
}

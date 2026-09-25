package domain

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       int    `json:"id" gorm:"primaryKey"`
	Username string `json:"username" binding:"required,min=3,max=32" gorm:"unique;not null"`
	Password string `json:"password,omitempty" binding:"required,min=6" gorm:"not null"`
}

func NewUser(
	id int,
	username string,
	password string,
) User {
	return User{
		ID:       id,
		Username: username,
		Password: password,
	}
}

const bcryptCost = 12

func (u *User) HashPassword() (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcryptCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}

	return string(hashedBytes), nil
}

func (u *User) CheckPassword(password string) bool {
	hash := u.Password
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) != nil
}

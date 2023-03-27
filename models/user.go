package models

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username  string     `gorm:"not null;unique" json:"username"`
	Password  string     `gorm:"not null" json:"password"`
	FullName  string     `gorm:"not null" json:"full_name"`
	Role      uint8      `gorm:"not null;size:1;default:0;comment:0 => staff, 1 => leader" json:"role"`
	UserScore UserScore  `json:"user_score"`
	UserTasks []TaskUser `json:"user_tasks"`
}

type UserResponse struct {
	ID        uint   `json:"id"`
	Username  string `json:"username"`
	FullName  string `json:"full_name"`
	Role      uint8  `json:"role"`
	UserScore UserScoreResponse
	UserTasks []UserTaskResponse
}

type UserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	FullName string `json:"full_name"`
}

type TokenRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (user *User) HashPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)

	if err != nil {
		return err
	}

	user.Password = string(bytes)

	return nil
}

func (user *User) CheckPassword(password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	if err != nil {
		return err
	}

	return nil
}

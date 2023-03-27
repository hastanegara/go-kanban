package config

import (
	"fmt"
	"project/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func DBConnect() {
	var err error
	USER := "root"
	PASS := ""
	HOST := "127.0.0.1"
	PORT := "3306"
	DBNAME := "kanban"
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", USER, PASS, HOST, PORT, DBNAME)
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	DB.AutoMigrate(&models.User{}, &models.Task{}, &models.TaskItem{}, &models.TaskDocument{}, &models.TaskUser{}, &models.UserScore{})

	bytes, _ := bcrypt.GenerateFromPassword([]byte("password"), 12)
	password := string(bytes)

	DB.Create(&models.User{
		Username: "leader",
		Password: password,
		FullName: "Leader Name",
		Role:     1,
	})
}

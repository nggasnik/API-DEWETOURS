package database

import (
	"erlangga/models"
	"erlangga/pkg/mysql"
	"fmt"
)

func RunMigration() {
	// menjalankan migration
	// err := mysql.DB.AutoMigrate(&models.Country{}, &models.Trip{}, &models.User{}, &models.Transaction{})
	err := mysql.DB.AutoMigrate(&models.Country{}, &models.Image{}, &models.Trip{}, &models.User{}, &models.Transaction{})
	if err != nil {
		fmt.Println(err)
		panic("Migration Failed")
	}

	fmt.Println("Migration Success")
}

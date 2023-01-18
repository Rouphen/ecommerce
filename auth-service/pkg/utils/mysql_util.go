package utils

import (
	"ecommerce/auth-service/pkg/domain"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func ConnectMySQLByGorm(url string) (*gorm.DB, error) {
	log.Println(url)
	db, err := gorm.Open(mysql.Open(url), &gorm.Config{})

	if err != nil {
		log.Fatalln(err)
		//return nil, err
	}

	db.AutoMigrate(&domain.User{})

	return db, nil
}

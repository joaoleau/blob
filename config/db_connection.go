package config

import (
	"log"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)
  

func ConnectToDB() {
	var err error
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=America/Sao_Paulo",
	cfg.DB.Host, cfg.DB.User, cfg.DB.Pass, cfg.DB.Database, cfg.DB.Port)
	DB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed to connect to database")
	}

}


package main

import (
	log "github.com/sirupsen/logrus"
	"payment/api/models"
	"payment/pkg/config"
	"payment/pkg/db"
)

func main() {
	configuration, err := config.LoadConfig("configs/config.yml")
	if err != nil {
		log.Fatal(err)
	}
	log.Info("Configuration loaded")

	dsn := db.LoadDSN(configuration)
	database, err := db.New(dsn)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to database")

	// Create enum types
	err = db.CreateEnumType(database, "discount_type", []string{"voucher", "charge"})
	if err != nil {
		log.Info(err)
	}
	err = db.CreateEnumType(database, "transaction_type", []string{"withdrawal", "deposit"})
	if err != nil {
		log.Info(err)
	}

	err = db.CreateEnumType(database, "transaction_status", []string{"pending", "completed", "failed"})
	if err != nil {
		log.Info(err)
	}

	// Migrate database models
	err = db.AutoMigrate(database,
		models.Wallet{},
		models.Transaction{})

	if err != nil {
		log.Fatal(err)
	}

	log.Info("Database models migrated")
}

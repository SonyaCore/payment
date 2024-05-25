package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"payment/api/models"
	"payment/pkg/config"
	"payment/pkg/db"
)

func main() {
	fmt.Println("Payment Service migration tool")
	logger := log.New()
	logger.SetFormatter(&log.TextFormatter{})

	configuration, err := config.LoadConfig("configs/config.yaml")
	if err != nil {
		logger.Fatal(err)
	}
	logger.Info("Configuration loaded")

	dsn := db.LoadDSN(configuration)
	database, err := db.New(dsn)
	if err != nil {
		logger.Fatal(err)
	}
	logger.Println("Connected to database")

	// Create enum types
	err = db.CreateEnumType(database, "discount_type", []string{"voucher", "charge"})
	if err != nil {
		logger.Info(err)
	}
	err = db.CreateEnumType(database, "transaction_type", []string{"withdrawal", "deposit"})
	if err != nil {
		logger.Info(err)
	}

	err = db.CreateEnumType(database, "transaction_status", []string{"pending", "completed", "failed"})
	if err != nil {
		logger.Info(err)
	}

	// Migrate database models
	err = db.AutoMigrate(database,
		models.Wallet{},
		models.Transaction{},
		models.Discount{},
		models.DiscountTransaction{},
	)

	if err != nil {
		logger.Fatal(err)
	}

	logger.Info("Database models migrated")
}

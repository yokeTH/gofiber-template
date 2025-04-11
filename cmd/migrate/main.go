package main

import (
	"fmt"
	"log"

	"github.com/yokeTH/gofiber-template/internal/config"
	"github.com/yokeTH/gofiber-template/internal/domain"
	"github.com/yokeTH/gofiber-template/pkg/db"
)

func main() {
	config := config.Load()

	db, err := db.New(config.PSQL)
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}

	if err := db.AutoMigrate(
		&domain.Book{},
		&domain.File{},
	); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	fmt.Println("Migration completed")
}

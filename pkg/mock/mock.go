package mock

import (
	"log"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func SetupMockDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("mock.db"), &gorm.Config{})
	return db, err
}

func CleanupMockDB() {
	err := os.Remove("mock.db")
	if err != nil {
		log.Fatalf("Failed to remove mock.db: %v", err)
	}
}

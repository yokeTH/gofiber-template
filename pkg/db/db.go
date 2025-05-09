package db

import (
	"fmt"
	"log"
	"math"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DBConfig struct {
	Host     string `env:"HOST,required"`
	Port     int    `env:"PORT,required"`
	User     string `env:"USER,required"`
	Password string `env:"PASSWORD,required"`
	DBName   string `env:"NAME,required"`
	SSLMode  string `env:"SSLMODE,required"`
}

func New(config DBConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s", config.Host, config.Port, config.User, config.Password, config.DBName, config.SSLMode)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		// Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatal("failed to connect database")
		return nil, err
	}

	return db, nil
}

// Pagination use with gorm.DB.Scopes() to calculate total pages and last page
//
// Usage Example:
//
//	db.Scopes(database.Paginate(domain.Book{}, limit, page, total, last)).Find(&books)
func Paginate(model any, limit, page, total, last *int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		var totalRows int64

		// create new gorm session for count row
		// idk why i have to do that
		// i just followed this https://stackoverflow.com/questions/72666748
		countDBSession := db.Session(&gorm.Session{Initialized: true})
		countDBSession.Model(model).Count(&totalRows)

		// db.Model(&domain.Book{}).Count(&totalRows)

		*total = int(totalRows)
		offset := (*page - 1) * *limit
		*last = int(math.Ceil(float64(totalRows) / float64(*limit)))

		return db.Offset(offset).Limit(*limit)
	}
}

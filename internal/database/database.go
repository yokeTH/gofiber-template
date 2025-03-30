package database

import (
	"fmt"
	"log"
	"math"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DBConfig struct {
	Host     string `env:"HOST"`
	Port     int    `env:"PORT"`
	User     string `env:"USER"`
	Password string `env:"PASSWORD"`
	DBName   string `env:"NAME"`
	SSLMode  string `env:"SSLMODE"`
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

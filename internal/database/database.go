package database

import (
	"errors"
	"math"

	"gorm.io/gorm"
)

type Database struct {
	*gorm.DB
}

func (db *Database) Paginate(value any, limit, page int, order string) (func(db *Database) *Database, int, int, error) {
	var totalRows int64
	if err := db.Model(value).Count(&totalRows).Error; err != nil {
		return nil, 0, 0, err
	}

	totalPages := int(math.Ceil(float64(totalRows) / float64(limit)))
	offset := (page - 1) * limit

	if page > totalPages {
		offset = (totalPages - 1) * limit
		return func(db *Database) *Database { return &Database{db.Offset(offset).Limit(limit).Order(order)} }, totalPages, int(totalRows), errors.New("page exceeded")
	}

	return func(db *Database) *Database { return &Database{db.Offset(offset).Limit(limit).Order(order)} }, totalPages, int(totalRows), nil

}

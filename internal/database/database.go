package database

import (
	"math"

	"gorm.io/gorm"
)

type Database struct {
	*gorm.DB
}

func (db *Database) Paginate(value any, limit int, page int, order string) (int, int, error) {
	var totalRows int64

	offset := (page - 1) * limit
	if err := db.Model(value).Count(&totalRows).Offset(offset).Limit(limit).Order(order).Find(value).Error; err != nil {
		return 0, 0, err
	}
	totalPages := int(math.Ceil(float64(totalRows) / float64(limit)))

	return totalPages, int(totalRows), nil

}

package domain

import (
	"gorm.io/gorm"
)

type Book struct {
	gorm.Model
	Title  string
	Author string
}

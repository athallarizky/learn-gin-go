package models

import "gorm.io/gorm"

type Article struct {
	gorm.Model
	Title string
	Slug  string `gorm:"uniqueIndex"`
	Desc  string `gorm:"type:text"`
}

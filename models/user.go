package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Email    string `gorm:"uniqueIndex"`
	Username string
	FullName string
	SocialId string
	Provider string
	Avatar   string
	Role     bool `gorm:"default:0"`
	Articles []Article
}

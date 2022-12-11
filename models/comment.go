package models

import "gorm.io/gorm"

type Comment struct {
	gorm.Model

	StrCommentDesc      string `gorm:"not null"`
	IntRate             int    `gorm:"not null"`
	IntCommentTimestamp int    `gorm:"not null"`
	IntUserId           uint   `gorm:"not null"`
	User                User   `gorm:"foreignKey:IntUserId"`
	IntRecipeId         uint   `gorm:"not null"`
	Recipe              Recipe `gorm:"foreignKey:IntRecipeId"`
}

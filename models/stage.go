package models

import "gorm.io/gorm"

type Stage struct {
	gorm.Model

	StrStageDesc string  `gorm:"not null"`
	IntRecipeId  uint    `gorm:"not null"`
	Recipe       Recipe  `gorm:"foreignKey:IntRecipeId"`
	StagePhotos  []Photo `gorm:"foreignKey:IntStageId"`
}

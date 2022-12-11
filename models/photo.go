package models

import "gorm.io/gorm"

type Photo struct {
	gorm.Model

	StrImage   string `gorm:"unique;not null"`
	IntStageId uint   `gorm:"not null"`
	Stage      Stage  `gorm:"foreignKey:IntStageId"`
}

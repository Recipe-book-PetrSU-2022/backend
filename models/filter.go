package models

import "gorm.io/gorm"

type Filter struct {
	gorm.Model

	StrFilterName string `gorm:"index;not null"`
}

package models

import "gorm.io/gorm"

type User struct {
	gorm.Model

	StrUserName     string    `gorm:"index;unique;not null"`
	StrUserPassword string    `gorm:"not null"`
	StrUserEmail    string    `gorm:"index;unique;not null"`
	IntUserRights   int       `gorm:"not null;default:0"`
	StrUserImage    string    `gorm:"not null;default:0"`
	UserRecipes     []Recipe  `gorm:"foreignKey:IntUserId"`
	UserComments    []Comment `gorm:"foreignKey:IntUserId"`
	UserFavorite    []Recipe  `gorm:"many2many:user_favorite_recipes"`
}

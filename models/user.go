package models

import "gorm.io/gorm"

type User struct {
	gorm.Model

	StrUserName     string    `gorm:"index;unique;not null"`
	StrUserPassword string    `gorm:"not null" json:"-"`
	StrUserEmail    string    `gorm:"index;unique;not null" json:"-"`
	IntUserRights   int       `gorm:"not null;default:0" json:"-"`
	StrUserImage    string    `gorm:"not null;default:0"`
	UserRecipes     []Recipe  `gorm:"foreignKey:IntUserId" json:"-"`
	UserComments    []Comment `gorm:"foreignKey:IntUserId" json:"-"`
	UserFavorite    []Recipe  `gorm:"many2many:user_favorite_recipes" json:"-"`
}

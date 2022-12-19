package models

import "gorm.io/gorm"

type Recipe struct {
	gorm.Model

	StrRecipeName        string             `gorm:"index;not null"`
	IntServings          int                `gorm:"not null;default:0"`
	IntTime              int                `gorm:"not null;default:0"`
	StrRecipeCountry     string             `gorm:"not null"`
	StrRecipeType        string             `gorm:"not null"`
	StrRecipeImage       string             `gorm:"not null"`
	BoolRecipeVisibility bool               `gorm:"not null"`
	IntUserId            uint               `gorm:"not null"`
	User                 User               `gorm:"foreignKey:IntUserId"`
	RecipeStages         []Stage            `gorm:"foreignKey:IntRecipeId"`
	RecipeComments       []Comment          `gorm:"foreignKey:IntRecipeId"`
	RecipeIngredients    []RecipeIngredient `gorm:"foreignKey:IntRecipeId"`
}

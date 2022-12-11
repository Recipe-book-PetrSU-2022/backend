package models

import "gorm.io/gorm"

type RecipeIngredient struct {
	gorm.Model

	IntCalories     int        `gorm:"not null;default:0"`
	IntRecipeId     uint       `gorm:"not null"`
	Recipe          Recipe     `gorm:"foreignKey:IntRecipeId"`
	IntIngredientId uint       `gorm:"not null"`
	Ingredient      Ingredient `gorm:"foreignKey:IntIngredientId"`
}

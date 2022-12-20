package models

import "gorm.io/gorm"

type RecipeIngredient struct {
	gorm.Model

	IntGrams        int        `gorm:"not null;default:0"`
	IntRecipeId     uint       `gorm:"not null;index:idx_recipe_ingr,unique"`
	Recipe          Recipe     `gorm:"foreignKey:IntRecipeId" json:"-"`
	IntIngredientId uint       `gorm:"not null;index:idx_recipe_ingr,unique"`
	Ingredient      Ingredient `gorm:"foreignKey:IntIngredientId"`
}

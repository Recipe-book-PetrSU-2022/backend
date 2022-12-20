package models

import "gorm.io/gorm"

type Ingredient struct {
	gorm.Model

	StrIngredientName string             `gorm:"unique;index;not null"`
	IntCalories       int                `gorm:"not null;default:0"`
	IntProteins       int                `gorm:"not null;default:0"`
	IntFats           int                `gorm:"not null;default:0"`
	IntCarbohydrates  int                `gorm:"not null;default:0"`
	RecipeIngredients []RecipeIngredient `gorm:"foreignKey:IntIngredientId" json:"-"`
}

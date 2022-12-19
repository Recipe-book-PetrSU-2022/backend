package main

// Структура обычного ответа
//
// Переменные структуры:
//   - Сообщение
type DefaultResponse struct {
	Message string `json:"message"` // Сообщение
}

// Структура ответа с токеном
//
// Переменные структуры:
//   - Сообщение
//   - Токен
type TokenResponse struct {
	Message string `json:"message"` // Сообщение
	Token   string `json:"token"`   // Токен
}

// Структура ответа с профилем пользователя
//
// Переменные структуры:
//   - Сообщение
//   - ID пользователя
//   - Никнейм
//   - Почта
//   - Фото профиля
type ProfileResponse struct {
	Message      string `json:"message"`  // Сообщение
	Id           uint   `json:"id"`       // ID пользователя
	Username     string `json:"username"` // Никнейм
	Email        string `json:"email"`    // Почта
	ProfilePhoto string `json:"photo"`    // Фото профиля
}

// Структура ответа с профилем пользователя
//
// Переменные структуры:
//   - Сообщение
//   - ID пользователя
//   - Никнейм
//   - Фото профиля
type UserResponse struct {
	Message      string `json:"message"`  // Сообщение
	Id           uint   `json:"id"`       // ID пользователя
	Username     string `json:"username"` // Никнейм
	ProfilePhoto string `json:"photo"`    // Фото профиля
}

type CommentResponse struct {
	Id        uint `json:"id"`
	UserId    uint `json:"user_id"`
	RecipieId uint `json:"recipie_id"`

	Description string `json:"description"`
	Rate        uint   `json:"rate"`
	Timestamp   int    `json:"timestamp"`
}

type RecipeIngredientResponse struct {
	Id           uint `json:"id"`
	RecipieId    uint `json:"recipie_id"`
	IngredientId uint `json:"ingredient_id"`

	Name          string `json:"name"`
	Grams         int    `json:"grams"`
	Calories      int    `json:"calories"`
	Proteins      int    `json:"proteins"`
	Fats          int    `json:"fats"`
	Carbohydrates int    `json:"carbohydrates"`
}

type StagePhotoResponose struct {
	Id    uint   `json:"id"`
	Image string `json:"name"`
}

type StageResponse struct {
	Id          uint                   `json:"id"`
	Description string                 `json:"description"`
	Photos      []*StagePhotoResponose `json:"photos"`
}

// Структура ответа с пустым рецептом
//
// Переменные структуры:
//   - Сообщение
type RecipeResponse struct {
	Message string `json:"message"` // Сообщение

	Id     uint `json:"id"`
	UserId uint `json:"user_id"`

	Name        string                      `json:"name"`
	Servings    int                         `json:"servings"`
	Time        int                         `json:"time"`
	Country     string                      `json:"country"`
	Type        string                      `json:"type"`
	Cover       string                      `json:"cover"`
	IsVisible   bool                        `json:"is_visible"`
	Stages      []*StageResponse            `json:"stages"`
	Comments    []*CommentResponse          `json:"comments"`
	Ingredients []*RecipeIngredientResponse `json:"ingredients"`
}

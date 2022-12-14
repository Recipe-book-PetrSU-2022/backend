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

// Структура ответа с пустым рецептом
//
// Переменные структуры:
//   - Сообщение
//   - ID рецепта
type RecipeResponse struct {
	Message string `json:"message"` // Сообщение
	Id      uint   `json:"id"`
}

// Структура ответа с этапом
//
// Переменные структуры:
//   - Сообщение
//   - ID этапа
type StageResponse struct {
	Message string `json:"message"` // Сообщение
	Id      uint   `json:"id"`
}

// Структура ответа с ингредиентом
//
// Переменные структуры:
//   - Сообщение
//   - ID этапа
type IngredientResponse struct {
	Message string `json:"message"` // Сообщение
	Id      uint   `json:"id"`
}

// Структура ответа с обложкой рецепта
//
// Переменные структуры:
//   - Сообщение
//   - Обложка рецепта
type CoverResponse struct {
	Message string `json:"message"` // Сообщение
	Cover   string `json:"cover"`
}

package main

import "golang.org/x/crypto/bcrypt"

// Функция для хэширования пароля
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// Функция для проверки на равенство между паролем и хэшом
func CheckPasswordHash(hash string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

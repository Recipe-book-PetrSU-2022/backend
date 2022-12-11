package main

type DefaultResponse struct {
	Message string `json:"message"`
}

type TokenResponse struct {
	Message string `json:"message"`
	Token   string `json:"token"`
}

type ProfileResponse struct {
	Message      string `json:"message"`
	Id           uint   `json:"id"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	ProfilePhoto string `json:"photo"`
}

type UserResponse struct {
	Message      string `json:"message"`
	Id           uint   `json:"id"`
	Username     string `json:"username"`
	ProfilePhoto string `json:"photo"`
}

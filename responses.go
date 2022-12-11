package main

type DefaultResponse struct {
	Message string `json:"message"`
}

type TokenResponse struct {
	Message string `json:"message"`
	Token   string `json:"token"`
}

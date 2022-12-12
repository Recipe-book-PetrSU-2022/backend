package main

import (
	"github.com/labstack/echo/v4"
)

type RecipeData struct {
	Name     string `json:"name"`
	Servings int    `json:"servings"`
	Time     int    `json:"time"`
	Country  string `json:"country"`
	Type     string `json:"type"`
	Image    string `json:"image"`
}

func (server *Server) CreateRecipeHandle(c echo.Context) error {
	return nil
}

func (server *Server) GetRecipeHandle(c echo.Context) error {
	return nil
}

func (server *Server) GetRecipesHandle(c echo.Context) error {
	return nil
}

func (server *Server) UpdateRecipeHandle(c echo.Context) error {
	return nil
}

func (server *Server) DeleteRecipeHandle(c echo.Context) error {
	return nil
}

func (server *Server) FindRecipesHandle(c echo.Context) error {
	return nil
}

func (server *Server) AddRecipeToFavoritesHandle(c echo.Context) error {
	return nil
}

func (server *Server) RemoveRecipeFromFavoritesHandle(c echo.Context) error {
	return nil
}

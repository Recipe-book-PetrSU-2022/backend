package main

import (
	"net/http"
	"os"
	"path"

	"github.com/labstack/echo/v4"
)

// Функция для скачивания файла
func (server *Server) DownloadFile(c echo.Context) error {
	// Получаем имя файла
	filename := c.Param("filename")
	if filename == "" {
		return c.JSON(http.StatusBadRequest, &DefaultResponse{
			Message: "Не указано имя файла",
		})
	}

	// Получаем путь к файлу
	filePath := path.Join(server.UploadsPath, filename)

	// Проверяем, существует ли файл
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return c.JSON(http.StatusNotFound, &DefaultResponse{
			Message: "Файл не найден",
		})
	}

	return c.File(filePath)
}

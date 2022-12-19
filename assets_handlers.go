package main

import (
	"net/http"
	"os"
	"path"

	"github.com/labstack/echo/v4"
)

func (server *Server) DownloadFile(c echo.Context) error {
	filename := c.Param("filename")
	if filename == "" {
		return c.JSON(http.StatusBadRequest, &DefaultResponse{
			Message: "Не указано имя файла",
		})
	}

	filePath := path.Join(server.UploadsPath, filename)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return c.JSON(http.StatusNotFound, &DefaultResponse{
			Message: "Файл не найден",
		})
	}

	return c.File(filePath)
}

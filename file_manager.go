package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path"
)

// Функция для создания директорий для хранения изображений
func (server *Server) CreateUploadDirs() error {
	fsErr := os.MkdirAll(server.UploadsPath, 0755)
	if fsErr != nil {
		return fsErr
	}
	return nil
}

// Функция для получения расширения файла
// Наследуется от Server, на вход принимает заголовок файла
func (server *Server) GetFileExtByMimetype(fileHeader *multipart.FileHeader) (string, error) {

	// Открываем файл
	file, err := fileHeader.Open()
	if err != nil {
		return "", err
	}

	// Берём первые 512 байт
	buff := make([]byte, 512)
	_, err = file.Read(buff)
	if err != nil {
		return "", err
	}

	// Выделяем тип файла из заголовка
	fileType := http.DetectContentType(buff)

	// Возвращаем тип файла
	switch fileType {
	case "image/png":
		return "png", nil
	case "image/jpeg":
		return "jpg", nil
	case "image/gif":
		return "gif", nil
	}

	// Если тип файла не PNG, JPG или GIF, то возвращаем ошибку
	return "", errors.New("разрешены только файлы png, jpg, gif")
}

// Функция для сохранения file с расширением ext в папку dst (см. выше)
func (server *Server) SaveFileWithExt(file multipart.File, ext string) (string, error) {
	// Создаем новое имя файла: случайная строка в 16 символов + расширение
	filename := fmt.Sprintf("%s.%s", RandomString(16), ext)
	// Создаем путь, по которому будет доступен файл
	filePath := path.Join(server.UploadsPath, filename)

	log.Printf("Saving file to %s", filePath)

	// Создаем файл
	dst, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	// Если не удается прочитать файл
	if _, err := io.Copy(dst, file); err != nil {
		// Т отправляем ошибку
		return "", err
	}

	// Иначе возвращаем имя файла
	return filename, nil
}

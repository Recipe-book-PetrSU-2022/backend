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

func (server *Server) CreateUploadDirs() error {
	fsErr := os.MkdirAll(server.UploadsPath, 0755)
	if fsErr != nil {
		return fsErr
	}
	return nil
}

func (server *Server) GetFileExtByMimetype(fileHeader *multipart.FileHeader) (string, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return "", err
	}

	buff := make([]byte, 512)
	_, err = file.Read(buff)
	if err != nil {
		return "", err
	}

	fileType := http.DetectContentType(buff)

	switch fileType {
	case "image/png":
		return "png", nil
	case "image/jpeg":
		return "jpg", nil
	case "image/gif":
		return "gif", nil
	}

	return "", errors.New("Разрешены только файлы png, jpg, gif")
}

// сохраняет file с расширением ext в папку dst (см. выше)
func (server *Server) SaveFileWithExt(file multipart.File, ext string) (string, error) {
	filename := fmt.Sprintf("%s.%s", RandomString(16), ext)
	filePath := path.Join(server.UploadsPath, filename)

	log.Printf("Saving file to %s", filePath)

	dst, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		return "", err
	}

	return filename, nil
}

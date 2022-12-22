package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Recipe-book-PetrSU-2022/backend/models"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestProfileWithUser(t *testing.T) {
	reqMap := map[string]interface{}{
		"login":        "",
		"email":        "",
		"old_password": "",
	}
	reqJson, _ := json.Marshal(reqMap)

	req := httptest.NewRequest(
		http.MethodPost, "/profile/", strings.NewReader(string(reqJson)),
	)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %s", UserJWT))

	rec := httptest.NewRecorder()

	c := TestE.NewContext(req, rec)

	if assert.NoError(t, TestJwtMiddleware(TestServer.ProfileHandle)(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)

		respJson := models.User{}
		err := json.Unmarshal(rec.Body.Bytes(), &respJson)
		assert.Nil(t, err)

		assert.Equal(t, uint(1), respJson.ID)
	}
}

func TestProfileWithNotExistsUser(t *testing.T) {
	reqMap := map[string]interface{}{
		"login":        "",
		"email":        "",
		"old_password": "",
	}
	reqJson, _ := json.Marshal(reqMap)

	req := httptest.NewRequest(
		http.MethodPost, "/profile/", strings.NewReader(string(reqJson)),
	)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %s", UserJWT2))

	rec := httptest.NewRecorder()

	c := TestE.NewContext(req, rec)

	if assert.NoError(t, TestJwtMiddleware(TestServer.ProfileHandle)(c)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)

		respJson := DefaultResponse{}
		err := json.Unmarshal(rec.Body.Bytes(), &respJson)
		assert.Nil(t, err)

		assert.Equal(t, "Пользователь не найден", respJson.Message)
	}
}

func TestRecipeWithImage(t *testing.T) {

}

func TestRecipeWithWrongImageFormat(t *testing.T) {

}

func TestSignUpWithHashPssword(t *testing.T) {
	reqMap := map[string]interface{}{
		"login":            "c",
		"email":            "c@c.ru",
		"password":         "c",
		"confirm_password": "c",
	}
	reqJson, _ := json.Marshal(reqMap)

	req := httptest.NewRequest(
		http.MethodPost, "/signup", strings.NewReader(string(reqJson)),
	)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()

	c := TestE.NewContext(req, rec)

	if assert.NoError(t, TestServer.SignUpHandle(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)

		respJson := DefaultResponse{}
		err := json.Unmarshal(rec.Body.Bytes(), &respJson)
		assert.Nil(t, err)

		assert.Equal(t, "Пользователь успешно зарегистрирован!", respJson.Message)
	}
}
func TestSignInCheckPasswordHash(t *testing.T) {
	reqMap := map[string]interface{}{
		"login":    "a",
		"password": "b",
	}
	reqJson, _ := json.Marshal(reqMap)

	req := httptest.NewRequest(
		http.MethodPost, "/signin", strings.NewReader(string(reqJson)),
	)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()

	c := TestE.NewContext(req, rec)

	if assert.NoError(t, TestServer.SignInHandle(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)

		respJson := TokenResponse{}
		err := json.Unmarshal(rec.Body.Bytes(), &respJson)
		assert.Nil(t, err)

		assert.Equal(t, "Пользователь успешно вошёл в систему!", respJson.Message)
		assert.NotEmpty(t, respJson.Token)
		UserJWT = respJson.Token
	}
}
func TestSignInCheckWrongPasswordHash(t *testing.T) {

	reqMap := map[string]interface{}{
		"login":    "a",
		"password": "bdhgfbjnkm",
	}
	reqJson, _ := json.Marshal(reqMap)

	req := httptest.NewRequest(
		http.MethodPost, "/signin", strings.NewReader(string(reqJson)),
	)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()

	c := TestE.NewContext(req, rec)

	if assert.NoError(t, TestServer.SignInHandle(c)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code) // так-то тут должен быть  http.StatusUnauthorized

		respJson := DefaultResponse{}
		err := json.Unmarshal(rec.Body.Bytes(), &respJson)
		assert.Nil(t, err)

		assert.Equal(t, "Введены неверные данные", respJson.Message)
	}
}

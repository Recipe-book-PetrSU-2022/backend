package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/Recipe-book-PetrSU-2022/backend/models"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	e          = echo.New()
	testServer = Server{
		E:    e,
		Host: "0.0.0.0",
		Port: 11111,
		// DBConnectionInfo: "file::memory:/test?cache=shared", // БД в оперативке
		TokenKey:    []byte("test"),
		UploadsPath: "/tmp/test/recipe_book_uploads",
	}
)

func TestMain(m *testing.M) {
	// err := testServer.ConnectDB()
	// if err != nil {
	// 	panic(err)
	// }

	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	testServer.DB = db

	err = testServer.DB.AutoMigrate(
		&models.User{},
		&models.Filter{},
		&models.Ingredient{},
		&models.Recipe{},
		&models.Stage{},
		&models.Comment{},
		&models.Photo{},
		&models.RecipeIngredient{},
	)
	if err != nil {
		panic(err)
	}

	os.Exit(m.Run())
}

func TestSignup(t *testing.T) {

	reqMap := map[string]interface{}{
		"login":            "a",
		"email":            "a@a.ru",
		"password":         "a",
		"confirm_password": "a",
	}
	reqJson, _ := json.Marshal(reqMap)

	req := httptest.NewRequest(
		http.MethodPost, "/signup", strings.NewReader(string(reqJson)),
	)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)

	if assert.NoError(t, testServer.SignUpHandle(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)

		respJson := DefaultResponse{}
		err := json.Unmarshal(rec.Body.Bytes(), &respJson)
		assert.Nil(t, err)

		assert.Equal(t, "Пользователь успешно зарегистрирован!", respJson.Message)
	}
}

func TestSignin(t *testing.T) {

	reqMap := map[string]interface{}{
		"login":    "a",
		"password": "a",
	}
	reqJson, _ := json.Marshal(reqMap)

	req := httptest.NewRequest(
		http.MethodPost, "/signin", strings.NewReader(string(reqJson)),
	)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)

	if assert.NoError(t, testServer.SignInHandle(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)

		respJson := TokenResponse{}
		err := json.Unmarshal(rec.Body.Bytes(), &respJson)
		assert.Nil(t, err)

		assert.Equal(t, "Пользователь успешно вошёл в систему!", respJson.Message)
		assert.NotEmpty(t, respJson.Token)
	}
}

func TestSigninInvalidCreds(t *testing.T) {

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

	c := e.NewContext(req, rec)

	if assert.NoError(t, testServer.SignInHandle(c)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code) // так-то тут должен быть  http.StatusUnauthorized

		respJson := DefaultResponse{}
		err := json.Unmarshal(rec.Body.Bytes(), &respJson)
		assert.Nil(t, err)

		assert.Equal(t, "Введены неверные данные", respJson.Message)
	}
}

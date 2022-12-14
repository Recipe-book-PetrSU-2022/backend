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

func TestSignupWithoutLogin(t *testing.T) {

	reqMap := map[string]interface{}{
		"login": "",
	}
	reqJson, _ := json.Marshal(reqMap)

	req := httptest.NewRequest(
		http.MethodPost, "/signup", strings.NewReader(string(reqJson)),
	)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()

	c := TestE.NewContext(req, rec)

	if assert.NoError(t, TestServer.SignUpHandle(c)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)

		respJson := DefaultResponse{}
		err := json.Unmarshal(rec.Body.Bytes(), &respJson)
		assert.Nil(t, err)

		assert.Equal(t, "Имя пользователя не может быть пустым", respJson.Message)
	}
}

func TestSignupWithoutEmail(t *testing.T) {

	reqMap := map[string]interface{}{
		"login": "a",
		"email": "",
	}
	reqJson, _ := json.Marshal(reqMap)

	req := httptest.NewRequest(
		http.MethodPost, "/signup", strings.NewReader(string(reqJson)),
	)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()

	c := TestE.NewContext(req, rec)

	if assert.NoError(t, TestServer.SignUpHandle(c)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)

		respJson := DefaultResponse{}
		err := json.Unmarshal(rec.Body.Bytes(), &respJson)
		assert.Nil(t, err)

		assert.Equal(t, "Почта пользователя не может быть пустой", respJson.Message)
	}
}

func TestSignupWithoutPassword(t *testing.T) {

	reqMap := map[string]interface{}{
		"login":    "a",
		"email":    "a@a.a",
		"password": "",
	}
	reqJson, _ := json.Marshal(reqMap)

	req := httptest.NewRequest(
		http.MethodPost, "/signup", strings.NewReader(string(reqJson)),
	)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()

	c := TestE.NewContext(req, rec)

	if assert.NoError(t, TestServer.SignUpHandle(c)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)

		respJson := DefaultResponse{}
		err := json.Unmarshal(rec.Body.Bytes(), &respJson)
		assert.Nil(t, err)

		assert.Equal(t, "Пароль пользователя не может быть пустым", respJson.Message)
	}
}

func TestSignupWithoutConfirmPassword(t *testing.T) {

	reqMap := map[string]interface{}{
		"login":            "a",
		"email":            "a@a.a",
		"password":         "a",
		"confirm_password": "",
	}
	reqJson, _ := json.Marshal(reqMap)

	req := httptest.NewRequest(
		http.MethodPost, "/signup", strings.NewReader(string(reqJson)),
	)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()

	c := TestE.NewContext(req, rec)

	if assert.NoError(t, TestServer.SignUpHandle(c)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)

		respJson := DefaultResponse{}
		err := json.Unmarshal(rec.Body.Bytes(), &respJson)
		assert.Nil(t, err)

		assert.Equal(t, "Пароли должны совпадать", respJson.Message)
	}
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

	c := TestE.NewContext(req, rec)

	if assert.NoError(t, TestServer.SignUpHandle(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)

		respJson := DefaultResponse{}
		err := json.Unmarshal(rec.Body.Bytes(), &respJson)
		assert.Nil(t, err)

		assert.Equal(t, "Пользователь успешно зарегистрирован!", respJson.Message)
	}
}

func TestSignup2(t *testing.T) {

	reqMap := map[string]interface{}{
		"login":            "b",
		"email":            "b@b.ru",
		"password":         "b",
		"confirm_password": "b",
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

func TestSignupUserExists(t *testing.T) {

	reqMap := map[string]interface{}{
		"login":            "a",
		"email":            "b@b.b",
		"password":         "b",
		"confirm_password": "b",
	}
	reqJson, _ := json.Marshal(reqMap)

	req := httptest.NewRequest(
		http.MethodPost, "/signup", strings.NewReader(string(reqJson)),
	)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()

	c := TestE.NewContext(req, rec)

	if assert.NoError(t, TestServer.SignUpHandle(c)) {
		assert.Equal(t, http.StatusInternalServerError, rec.Code)

		respJson := DefaultResponse{}
		err := json.Unmarshal(rec.Body.Bytes(), &respJson)
		assert.Nil(t, err)

		assert.Equal(t, "Не получилось создать пользователя", respJson.Message)
	}
}

func TestSigninWithoutLogin(t *testing.T) {

	reqMap := map[string]interface{}{
		"login": "",
	}
	reqJson, _ := json.Marshal(reqMap)

	req := httptest.NewRequest(
		http.MethodPost, "/signin", strings.NewReader(string(reqJson)),
	)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()

	c := TestE.NewContext(req, rec)

	if assert.NoError(t, TestServer.SignInHandle(c)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)

		respJson := DefaultResponse{}
		err := json.Unmarshal(rec.Body.Bytes(), &respJson)
		assert.Nil(t, err)

		assert.Equal(t, "Имя пользователя не может быть пустым", respJson.Message)
	}
}

func TestSigninWithoutPassword(t *testing.T) {

	reqMap := map[string]interface{}{
		"login":    "a",
		"password": "",
	}
	reqJson, _ := json.Marshal(reqMap)

	req := httptest.NewRequest(
		http.MethodPost, "/signin", strings.NewReader(string(reqJson)),
	)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()

	c := TestE.NewContext(req, rec)

	if assert.NoError(t, TestServer.SignInHandle(c)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)

		respJson := DefaultResponse{}
		err := json.Unmarshal(rec.Body.Bytes(), &respJson)
		assert.Nil(t, err)

		assert.Equal(t, "Пароль пользователя не может быть пустым", respJson.Message)
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

func TestSignin2(t *testing.T) {

	reqMap := map[string]interface{}{
		"login":    "b",
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
		UserJWT2 = respJson.Token
	}
}

func TestSigninWithNotExistsUser(t *testing.T) {

	reqMap := map[string]interface{}{
		"login":    "hehehehe",
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
		assert.Equal(t, http.StatusBadRequest, rec.Code) // так-то тут должен быть  http.StatusUnauthorized

		respJson := DefaultResponse{}
		err := json.Unmarshal(rec.Body.Bytes(), &respJson)
		assert.Nil(t, err)

		assert.Equal(t, "Пользователь не найден", respJson.Message)
	}
}

func TestSigninWithWrongPassword(t *testing.T) {

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
		assert.Equal(t, http.StatusBadRequest, rec.Code) // так-то тут должен быть  http.StatusUnauthorized

		respJson := DefaultResponse{}
		err := json.Unmarshal(rec.Body.Bytes(), &respJson)
		assert.Nil(t, err)

		assert.Equal(t, "Введены неверные данные", respJson.Message)
	}
}

func TestProfile(t *testing.T) {

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

func TestChangeUserInfoNothing(t *testing.T) {

	reqMap := map[string]interface{}{
		"login":        "",
		"email":        "",
		"old_password": "",
	}
	reqJson, _ := json.Marshal(reqMap)

	req := httptest.NewRequest(
		http.MethodPost, "/profile/update", strings.NewReader(string(reqJson)),
	)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %s", UserJWT))

	rec := httptest.NewRecorder()

	c := TestE.NewContext(req, rec)

	if assert.NoError(t, TestJwtMiddleware(TestServer.ChangeProfileHandle)(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)

		respJson := DefaultResponse{}
		err := json.Unmarshal(rec.Body.Bytes(), &respJson)
		assert.Nil(t, err)

		assert.Equal(t, "Нечего изменять", respJson.Message)
	}
}

func TestChangeUserInfoEmail(t *testing.T) {

	reqMap := map[string]interface{}{
		"email": "new_email@a.ru",
	}
	reqJson, _ := json.Marshal(reqMap)

	req := httptest.NewRequest(
		http.MethodPost, "/profile/update", strings.NewReader(string(reqJson)),
	)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %s", UserJWT))

	rec := httptest.NewRecorder()

	c := TestE.NewContext(req, rec)

	if assert.NoError(t, TestJwtMiddleware(TestServer.ChangeProfileHandle)(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)

		respJson := DefaultResponse{}
		err := json.Unmarshal(rec.Body.Bytes(), &respJson)
		assert.Nil(t, err)

		assert.Equal(t, "Пользователь успешно изменил свои данные!", respJson.Message)

		var user models.User
		TestServer.DB.First(&user, "id = ?", 1)
		assert.Equal(t, "new_email@a.ru", user.StrUserEmail)
	}
}

func TestChangeUserInfoWithoutOldPassword(t *testing.T) {

	reqMap := map[string]interface{}{
		"password": "b",
	}
	reqJson, _ := json.Marshal(reqMap)

	req := httptest.NewRequest(
		http.MethodPost, "/profile/update", strings.NewReader(string(reqJson)),
	)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %s", UserJWT))

	rec := httptest.NewRecorder()

	c := TestE.NewContext(req, rec)

	if assert.NoError(t, TestJwtMiddleware(TestServer.ChangeProfileHandle)(c)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)

		respJson := DefaultResponse{}
		err := json.Unmarshal(rec.Body.Bytes(), &respJson)
		assert.Nil(t, err)

		assert.Equal(t, "Вы не ввели текущий пароль", respJson.Message)
	}
}

func TestChangeUserInfoWithWrongConfirmPassword(t *testing.T) {

	reqMap := map[string]interface{}{
		"old_password":     "a",
		"password":         "b",
		"confirm_password": "c",
	}
	reqJson, _ := json.Marshal(reqMap)

	req := httptest.NewRequest(
		http.MethodPost, "/profile/update", strings.NewReader(string(reqJson)),
	)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %s", UserJWT))

	rec := httptest.NewRecorder()

	c := TestE.NewContext(req, rec)

	if assert.NoError(t, TestJwtMiddleware(TestServer.ChangeProfileHandle)(c)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)

		respJson := DefaultResponse{}
		err := json.Unmarshal(rec.Body.Bytes(), &respJson)
		assert.Nil(t, err)

		assert.Equal(t, "Новые пароли должны совпадать", respJson.Message)
	}
}

func TestChangeUserInfoWithWrongOldPassword(t *testing.T) {

	reqMap := map[string]interface{}{
		"old_password":     "b",
		"password":         "b",
		"confirm_password": "b",
	}
	reqJson, _ := json.Marshal(reqMap)

	req := httptest.NewRequest(
		http.MethodPost, "/profile/update", strings.NewReader(string(reqJson)),
	)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %s", UserJWT))

	rec := httptest.NewRecorder()

	c := TestE.NewContext(req, rec)

	if assert.NoError(t, TestJwtMiddleware(TestServer.ChangeProfileHandle)(c)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)

		respJson := DefaultResponse{}
		err := json.Unmarshal(rec.Body.Bytes(), &respJson)
		assert.Nil(t, err)

		assert.Equal(t, "Введен неверный пароль", respJson.Message)
	}
}

func TestChangeUserInfoPassword(t *testing.T) {

	reqMap := map[string]interface{}{
		"old_password":     "a",
		"password":         "b",
		"confirm_password": "b",
	}
	reqJson, _ := json.Marshal(reqMap)

	req := httptest.NewRequest(
		http.MethodPost, "/profile/update", strings.NewReader(string(reqJson)),
	)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %s", UserJWT))

	rec := httptest.NewRecorder()

	c := TestE.NewContext(req, rec)

	if assert.NoError(t, TestJwtMiddleware(TestServer.ChangeProfileHandle)(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)

		respJson := DefaultResponse{}
		err := json.Unmarshal(rec.Body.Bytes(), &respJson)
		assert.Nil(t, err)

		assert.Equal(t, "Пользователь успешно изменил свои данные!", respJson.Message)
	}
}

func TestCreateEmptyRecipe(t *testing.T) {

	reqMap := map[string]interface{}{
		"old_password":     "a",
		"password":         "b",
		"confirm_password": "b",
	}
	reqJson, _ := json.Marshal(reqMap)

	req := httptest.NewRequest(
		http.MethodPost, "/recipe/add", strings.NewReader(string(reqJson)),
	)
	// req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %s", UserJWT))

	rec := httptest.NewRecorder()

	c := TestE.NewContext(req, rec)

	if assert.NoError(t, TestJwtMiddleware(TestServer.CreateEmptyRecipeHandle)(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)

		respJson := RecipeResponse{}
		err := json.Unmarshal(rec.Body.Bytes(), &respJson)
		assert.Nil(t, err)
		id := respJson.Id

		var recipe models.Recipe
		err = TestServer.DB.First(&recipe, "id = ?", id).Error
		assert.Nil(t, err)

		assert.Equal(t, "Создан новый рецепт", respJson.Message)
	}
}

func TestUpdateNotExistsRecipe(t *testing.T) {

	reqMap := map[string]interface{}{
		"old_password":     "a",
		"password":         "b",
		"confirm_password": "b",
	}
	reqJson, _ := json.Marshal(reqMap)

	req := httptest.NewRequest(
		http.MethodPost, "/my-recipe/complete/1000", strings.NewReader(string(reqJson)),
	)
	// req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %s", UserJWT))

	rec := httptest.NewRecorder()

	c := TestE.NewContext(req, rec)
	c.SetPath("/my-recipe/complete/1000")
	c.SetParamNames("id")
	c.SetParamValues("1000")

	if assert.NoError(t, TestJwtMiddleware(TestServer.UpdateRecipeHandle)(c)) {
		assert.Equal(t, http.StatusInternalServerError, rec.Code)

		respJson := RecipeResponse{}
		err := json.Unmarshal(rec.Body.Bytes(), &respJson)
		assert.Nil(t, err)
		assert.Equal(t, "Не удалось найти рецепт", respJson.Message)
	}
}

func TestUpdateRecipeWithoutName(t *testing.T) {

	reqMap := map[string]interface{}{
		"name": "",
	}
	reqJson, _ := json.Marshal(reqMap)

	req := httptest.NewRequest(
		http.MethodPost, "/recipe/complete/1", strings.NewReader(string(reqJson)),
	)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %s", UserJWT))

	rec := httptest.NewRecorder()

	c := TestE.NewContext(req, rec)
	c.SetPath("/my-recipe/complete/1")
	c.SetParamNames("id")
	c.SetParamValues("1")

	if assert.NoError(t, TestJwtMiddleware(TestServer.UpdateRecipeHandle)(c)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)

		respJson := RecipeResponse{}
		err := json.Unmarshal(rec.Body.Bytes(), &respJson)
		assert.Nil(t, err)
		assert.Equal(t, "Название рецепта не может быть пустым", respJson.Message)
	}
}

func TestUpdateRecipe(t *testing.T) {

	reqMap := map[string]interface{}{
		"name": "b",
	}
	reqJson, _ := json.Marshal(reqMap)

	req := httptest.NewRequest(
		http.MethodPost, "/recipe/add", strings.NewReader(string(reqJson)),
	)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %s", UserJWT))

	rec := httptest.NewRecorder()

	c := TestE.NewContext(req, rec)
	c.SetPath("/my-recipe/complete/1")
	c.SetParamNames("id")
	c.SetParamValues("1")

	if assert.NoError(t, TestJwtMiddleware(TestServer.UpdateRecipeHandle)(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)

		respJson := RecipeResponse{}
		err := json.Unmarshal(rec.Body.Bytes(), &respJson)
		assert.Nil(t, err)
		id := 1

		var recipe models.Recipe
		err = TestServer.DB.First(&recipe, "id = ?", id).Error
		assert.Equal(t, "b", recipe.StrRecipeName)

		assert.Equal(t, "Рецепт обновлен", respJson.Message)
	}
}

func TestUpdateRecipeByAnother(t *testing.T) {

	reqMap := map[string]interface{}{
		"name": "b",
	}
	reqJson, _ := json.Marshal(reqMap)

	req := httptest.NewRequest(
		http.MethodPost, "/recipe/add", strings.NewReader(string(reqJson)),
	)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %s", UserJWT2))

	rec := httptest.NewRecorder()

	c := TestE.NewContext(req, rec)
	c.SetPath("/my-recipe/complete/1")
	c.SetParamNames("id")
	c.SetParamValues("1")

	if assert.NoError(t, TestJwtMiddleware(TestServer.UpdateRecipeHandle)(c)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)

		respJson := RecipeResponse{}
		err := json.Unmarshal(rec.Body.Bytes(), &respJson)
		assert.Nil(t, err)
		assert.Equal(t, "Рецепт принадлежит другому пользователю", respJson.Message)
	}
}

func TestGetNotExistsRecipe(t *testing.T) {

	reqMap := map[string]interface{}{
		"name": "b",
	}
	reqJson, _ := json.Marshal(reqMap)

	req := httptest.NewRequest(
		http.MethodGet, "/recipe/1000", strings.NewReader(string(reqJson)),
	)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %s", UserJWT))

	rec := httptest.NewRecorder()

	c := TestE.NewContext(req, rec)
	c.SetPath("/recipe/1000")
	c.SetParamNames("id")
	c.SetParamValues("1000")

	if assert.NoError(t, TestJwtMiddleware(TestServer.GetRecipeHandle)(c)) {
		assert.Equal(t, http.StatusInternalServerError, rec.Code)

		respJson := DefaultResponse{}
		err := json.Unmarshal(rec.Body.Bytes(), &respJson)
		assert.Nil(t, err)

		assert.Equal(t, "Не удалось найти рецепт", respJson.Message)
	}
}

func TestChangeVisibilitySetVisibleRecipe(t *testing.T) {

	reqMap := map[string]interface{}{
		"visible": true,
	}
	reqJson, _ := json.Marshal(reqMap)

	req := httptest.NewRequest(
		http.MethodGet, "/visible/1", strings.NewReader(string(reqJson)),
	)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %s", UserJWT))

	rec := httptest.NewRecorder()

	c := TestE.NewContext(req, rec)
	c.SetPath("/recipe/1")
	c.SetParamNames("id")
	c.SetParamValues("1")

	if assert.NoError(t, TestJwtMiddleware(TestServer.ChangeVisibilityRecipeHandle)(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)

		respJson := DefaultResponse{}
		err := json.Unmarshal(rec.Body.Bytes(), &respJson)
		assert.Nil(t, err)

		assert.Equal(t, "Рецепт обновлен", respJson.Message)
	}
}

func TestChangeVisibilitySetInvisibleRecipe(t *testing.T) {

	reqMap := map[string]interface{}{
		"visible": false,
	}
	reqJson, _ := json.Marshal(reqMap)

	req := httptest.NewRequest(
		http.MethodGet, "/visible/1", strings.NewReader(string(reqJson)),
	)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %s", UserJWT))

	rec := httptest.NewRecorder()

	c := TestE.NewContext(req, rec)
	c.SetPath("/recipe/1")
	c.SetParamNames("id")
	c.SetParamValues("1")

	if assert.NoError(t, TestJwtMiddleware(TestServer.ChangeVisibilityRecipeHandle)(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)

		respJson := DefaultResponse{}
		err := json.Unmarshal(rec.Body.Bytes(), &respJson)
		assert.Nil(t, err)

		assert.Equal(t, "Рецепт обновлен", respJson.Message)
	}
}

func TestGetHiddenRecipe(t *testing.T) {

	reqMap := map[string]interface{}{
		"name": "b",
	}
	reqJson, _ := json.Marshal(reqMap)

	req := httptest.NewRequest(
		http.MethodGet, "/recipe/1", strings.NewReader(string(reqJson)),
	)
	// req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %s", UserJWT))

	rec := httptest.NewRecorder()

	c := TestE.NewContext(req, rec)
	c.SetPath("/recipe/1")
	c.SetParamNames("id")
	c.SetParamValues("1")

	if assert.NoError(t, TestJwtMiddleware(TestServer.GetRecipeHandle)(c)) {
		assert.Equal(t, http.StatusInternalServerError, rec.Code)

		respJson := DefaultResponse{}
		err := json.Unmarshal(rec.Body.Bytes(), &respJson)
		assert.Nil(t, err)

		assert.Equal(t, "Не удалось найти рецепт", respJson.Message)
	}

	// TestChangeVisibilityRecipe(t)

	t.Run("Revert visiblilty", TestChangeVisibilitySetVisibleRecipe)
}

func TestGetRecipe(t *testing.T) {

	reqMap := map[string]interface{}{
		"name": "b",
	}
	reqJson, _ := json.Marshal(reqMap)

	req := httptest.NewRequest(
		http.MethodGet, "/recipe/1", strings.NewReader(string(reqJson)),
	)
	// req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %s", UserJWT))

	rec := httptest.NewRecorder()

	c := TestE.NewContext(req, rec)
	c.SetPath("/recipe/1")
	c.SetParamNames("id")
	c.SetParamValues("1")

	if assert.NoError(t, TestJwtMiddleware(TestServer.GetRecipeHandle)(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)

		respJson := models.Recipe{}
		err := json.Unmarshal(rec.Body.Bytes(), &respJson)
		assert.Nil(t, err)
		var recipe models.Recipe
		err = TestServer.DB.First(&recipe, "id = ?", 1).Error
		assert.Nil(t, err)

		assert.Equal(t, recipe.ID, respJson.ID)
		assert.Equal(t, recipe.StrRecipeName, respJson.StrRecipeName)
	}
}

func TestGetPersonalNotExistsRecipe(t *testing.T) {

	reqMap := map[string]interface{}{
		"name": "b",
	}
	reqJson, _ := json.Marshal(reqMap)

	req := httptest.NewRequest(
		http.MethodGet, "/recipe/1000", strings.NewReader(string(reqJson)),
	)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %s", UserJWT))

	rec := httptest.NewRecorder()

	c := TestE.NewContext(req, rec)
	c.SetPath("/recipe/1000")
	c.SetParamNames("id")
	c.SetParamValues("1000")

	if assert.NoError(t, TestJwtMiddleware(TestServer.GetMyRecipeHandle)(c)) {
		assert.Equal(t, http.StatusInternalServerError, rec.Code)

		respJson := DefaultResponse{}
		err := json.Unmarshal(rec.Body.Bytes(), &respJson)
		assert.Nil(t, err)

		assert.Equal(t, "Не удалось найти рецепт", respJson.Message)
	}
}

func TestGetPersonalRecipeByAnother(t *testing.T) {
	reqMap := map[string]interface{}{
		"name": "b",
	}
	reqJson, _ := json.Marshal(reqMap)

	req := httptest.NewRequest(
		http.MethodGet, "/recipe/1", strings.NewReader(string(reqJson)),
	)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %s", UserJWT2))

	rec := httptest.NewRecorder()

	c := TestE.NewContext(req, rec)
	c.SetPath("/recipe/1")
	c.SetParamNames("id")
	c.SetParamValues("1")

	if assert.NoError(t, TestJwtMiddleware(TestServer.GetMyRecipeHandle)(c)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)

		respJson := DefaultResponse{}
		err := json.Unmarshal(rec.Body.Bytes(), &respJson)
		assert.Nil(t, err)

		assert.Equal(t, "Рецепт принадлежит другому пользователю", respJson.Message)
	}
}

func TestGetPersonalRecipe(t *testing.T) {

	reqMap := map[string]interface{}{
		"name": "b",
	}
	reqJson, _ := json.Marshal(reqMap)

	req := httptest.NewRequest(
		http.MethodGet, "/recipe/1", strings.NewReader(string(reqJson)),
	)
	// req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %s", UserJWT))

	rec := httptest.NewRecorder()

	c := TestE.NewContext(req, rec)
	c.SetPath("/recipe/1")
	c.SetParamNames("id")
	c.SetParamValues("1")

	if assert.NoError(t, TestJwtMiddleware(TestServer.GetMyRecipeHandle)(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)

		respJson := models.Recipe{}
		err := json.Unmarshal(rec.Body.Bytes(), &respJson)
		assert.Nil(t, err)
		var recipe models.Recipe
		err = TestServer.DB.First(&recipe, "id = ?", 1).Error
		assert.Nil(t, err)

		assert.Equal(t, recipe.ID, respJson.ID)
		assert.Equal(t, recipe.StrRecipeName, respJson.StrRecipeName)
	}
}

func TestGetAllRecipies(t *testing.T) {

	reqMap := map[string]interface{}{
		"name": "b",
	}
	reqJson, _ := json.Marshal(reqMap)

	req := httptest.NewRequest(
		http.MethodGet, "/recipe/all", strings.NewReader(string(reqJson)),
	)
	// req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %s", UserJWT))

	rec := httptest.NewRecorder()

	c := TestE.NewContext(req, rec)

	if assert.NoError(t, TestJwtMiddleware(TestServer.GetRecipesHandle)(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)

		respJson := []models.Recipe{}
		err := json.Unmarshal(rec.Body.Bytes(), &respJson)
		assert.Nil(t, err)
		var recipies []models.Recipe
		err = TestServer.DB.Find(&recipies).Error
		assert.Nil(t, err)

		assert.Equal(t, recipies[0].ID, respJson[0].ID)
		assert.Equal(t, len(recipies), len(respJson))
	}
}

func TestGetAllPersonalRecipies(t *testing.T) {

	reqMap := map[string]interface{}{
		"name": "b",
	}
	reqJson, _ := json.Marshal(reqMap)

	req := httptest.NewRequest(
		http.MethodGet, "/my-recipe/all", strings.NewReader(string(reqJson)),
	)
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %s", UserJWT))

	rec := httptest.NewRecorder()

	c := TestE.NewContext(req, rec)

	if assert.NoError(t, TestJwtMiddleware(TestServer.GetMyRecipesHandle)(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)

		respJson := []models.Recipe{}
		err := json.Unmarshal(rec.Body.Bytes(), &respJson)
		assert.Nil(t, err)
		var recipies []models.Recipe
		err = TestServer.DB.Find(&recipies, "int_user_id = ?", 1).Error
		assert.Nil(t, err)

		assert.Equal(t, recipies[0].ID, respJson[0].ID)
		assert.Equal(t, len(recipies), len(respJson))
	}
}

func TestDeleteNotExistsRecipe(t *testing.T) {

	reqMap := map[string]interface{}{
		"name": "b",
	}
	reqJson, _ := json.Marshal(reqMap)

	req := httptest.NewRequest(
		http.MethodGet, "/my-recipe/delete/100", strings.NewReader(string(reqJson)),
	)
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %s", UserJWT))

	rec := httptest.NewRecorder()

	c := TestE.NewContext(req, rec)
	c.SetPath("/my-recipe/delete/100")
	c.SetParamNames("id")
	c.SetParamValues("100")

	if assert.NoError(t, TestJwtMiddleware(TestServer.DeleteRecipeHandle)(c)) {
		assert.Equal(t, http.StatusInternalServerError, rec.Code)

		respJson := DefaultResponse{}
		err := json.Unmarshal(rec.Body.Bytes(), &respJson)
		assert.Nil(t, err)
		assert.Equal(t, "Не удалось найти рецепт", respJson.Message)
	}
}

func TestDeleteRecipeByAnother(t *testing.T) {

	reqMap := map[string]interface{}{
		"name": "b",
	}
	reqJson, _ := json.Marshal(reqMap)

	req := httptest.NewRequest(
		http.MethodGet, "/my-recipe/delete/1", strings.NewReader(string(reqJson)),
	)
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %s", UserJWT2))

	rec := httptest.NewRecorder()

	c := TestE.NewContext(req, rec)
	c.SetPath("/my-recipe/delete/1")
	c.SetParamNames("id")
	c.SetParamValues("1")

	if assert.NoError(t, TestJwtMiddleware(TestServer.DeleteRecipeHandle)(c)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)

		respJson := DefaultResponse{}
		err := json.Unmarshal(rec.Body.Bytes(), &respJson)
		assert.Nil(t, err)
		assert.Equal(t, "Рецепт принадлежит другому пользователю", respJson.Message)
	}
}

func TestFindRecipeByEmptySearch(t *testing.T) {

	reqMap := map[string]interface{}{
		"text": "",
	}
	reqJson, _ := json.Marshal(reqMap)

	req := httptest.NewRequest(
		http.MethodGet, "/recipe/find", strings.NewReader(string(reqJson)),
	)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %s", UserJWT))

	rec := httptest.NewRecorder()

	c := TestE.NewContext(req, rec)

	if assert.NoError(t, TestJwtMiddleware(TestServer.FindRecipesHandle)(c)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)

		respJson := DefaultResponse{}
		err := json.Unmarshal(rec.Body.Bytes(), &respJson)
		assert.Nil(t, err)

		assert.Equal(t, "Пустая строка поиска", respJson.Message)
	}
}

func TestFindRecipeByWeirdSearch(t *testing.T) {

	reqMap := map[string]interface{}{
		"text": "gfhjkl;lkjhgfv",
	}
	reqJson, _ := json.Marshal(reqMap)

	req := httptest.NewRequest(
		http.MethodGet, "/recipe/find", strings.NewReader(string(reqJson)),
	)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %s", UserJWT))

	rec := httptest.NewRecorder()

	c := TestE.NewContext(req, rec)

	if assert.NoError(t, TestJwtMiddleware(TestServer.FindRecipesHandle)(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)

		respJson := []models.Recipe{}
		err := json.Unmarshal(rec.Body.Bytes(), &respJson)
		assert.Nil(t, err)

		assert.Equal(t, 0, len(respJson))
	}
}

func TestFindRecipe(t *testing.T) {

	reqMap := map[string]interface{}{
		"text": "b",
	}
	reqJson, _ := json.Marshal(reqMap)

	req := httptest.NewRequest(
		http.MethodGet, "/recipe/find", strings.NewReader(string(reqJson)),
	)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %s", UserJWT))

	rec := httptest.NewRecorder()

	c := TestE.NewContext(req, rec)

	if assert.NoError(t, TestJwtMiddleware(TestServer.FindRecipesHandle)(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)

		respJson := []models.Recipe{}
		err := json.Unmarshal(rec.Body.Bytes(), &respJson)
		assert.Nil(t, err)

		assert.Equal(t, 1, len(respJson))
	}
}

func TestAddNotExistsRecipeToFavorites(t *testing.T) {

	reqMap := map[string]interface{}{
		"text": "b",
	}
	reqJson, _ := json.Marshal(reqMap)

	req := httptest.NewRequest(
		http.MethodGet, "/recipe/favorite/100", strings.NewReader(string(reqJson)),
	)
	// req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %s", UserJWT))

	rec := httptest.NewRecorder()

	c := TestE.NewContext(req, rec)
	c.SetPath("/my-recipe/favorite/100")
	c.SetParamNames("id")
	c.SetParamValues("100")

	if assert.NoError(t, TestJwtMiddleware(TestServer.AddRecipeToFavoritesHandle)(c)) {
		assert.Equal(t, http.StatusInternalServerError, rec.Code)

		respJson := DefaultResponse{}
		err := json.Unmarshal(rec.Body.Bytes(), &respJson)
		assert.Nil(t, err)

		assert.Equal(t, "Не удалось найти рецепт", respJson.Message)
	}
}

func TestAddRecipeToFavorites(t *testing.T) {

	reqMap := map[string]interface{}{
		"text": "b",
	}
	reqJson, _ := json.Marshal(reqMap)

	req := httptest.NewRequest(
		http.MethodGet, "/recipe/favorite/1", strings.NewReader(string(reqJson)),
	)
	// req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %s", UserJWT))

	rec := httptest.NewRecorder()

	c := TestE.NewContext(req, rec)
	c.SetPath("/my-recipe/favorite/1")
	c.SetParamNames("id")
	c.SetParamValues("1")

	if assert.NoError(t, TestJwtMiddleware(TestServer.AddRecipeToFavoritesHandle)(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)

		respJson := DefaultResponse{}
		err := json.Unmarshal(rec.Body.Bytes(), &respJson)
		assert.Nil(t, err)

		assert.Equal(t, "Ок", respJson.Message)
	}
}

func TestChangeVisibilitySetInvisibleRecipeByAnother(t *testing.T) {

	reqMap := map[string]interface{}{
		"visible": false,
	}
	reqJson, _ := json.Marshal(reqMap)

	req := httptest.NewRequest(
		http.MethodGet, "/visible/1", strings.NewReader(string(reqJson)),
	)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %s", UserJWT2))

	rec := httptest.NewRecorder()

	c := TestE.NewContext(req, rec)
	c.SetPath("/recipe/1")
	c.SetParamNames("id")
	c.SetParamValues("1")

	if assert.NoError(t, TestJwtMiddleware(TestServer.ChangeVisibilityRecipeHandle)(c)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)

		respJson := DefaultResponse{}
		err := json.Unmarshal(rec.Body.Bytes(), &respJson)
		assert.Nil(t, err)

		assert.Equal(t, "Рецепт принадлежит другому пользователю", respJson.Message)
	}
}

func TestGetNotExistsComment(t *testing.T) {

	reqMap := map[string]interface{}{
		"visible": false,
	}
	reqJson, _ := json.Marshal(reqMap)

	req := httptest.NewRequest(
		http.MethodGet, "/recipe/1/comment/1", strings.NewReader(string(reqJson)),
	)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %s", UserJWT2))

	rec := httptest.NewRecorder()

	c := TestE.NewContext(req, rec)
	c.SetPath("/recipe/1/comment/1")
	c.SetParamNames("recipe_id", "comment_id")
	c.SetParamValues("1", "1")

	if assert.NoError(t, TestJwtMiddleware(TestServer.GetCommentHandle)(c)) {
		assert.Equal(t, http.StatusNotFound, rec.Code)

		respJson := DefaultResponse{}
		err := json.Unmarshal(rec.Body.Bytes(), &respJson)
		assert.Nil(t, err)

		assert.Equal(t, "Комментарий не найден", respJson.Message)
	}
}

func TestCreateCommentWithNegativeRate(t *testing.T) {

	reqMap := map[string]interface{}{
		"text": "",
		"rate": -1,
	}
	reqJson, _ := json.Marshal(reqMap)

	req := httptest.NewRequest(
		http.MethodGet, "/recipe/1/comment/1", strings.NewReader(string(reqJson)),
	)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %s", UserJWT2))

	rec := httptest.NewRecorder()

	c := TestE.NewContext(req, rec)
	c.SetPath("/recipe/1/comment/add")
	c.SetParamNames("recipe_id")
	c.SetParamValues("1")

	if assert.NoError(t, TestJwtMiddleware(TestServer.CreateCommentHandle)(c)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)

		respJson := DefaultResponse{}
		err := json.Unmarshal(rec.Body.Bytes(), &respJson)
		assert.Nil(t, err)

		assert.Equal(t, "Оценка пользователя должна быть от 0 до 5", respJson.Message)
	}
}

func TestCreateCommentWithBiggerRate(t *testing.T) {

	reqMap := map[string]interface{}{
		"text": "",
		"rate": 6,
	}
	reqJson, _ := json.Marshal(reqMap)

	req := httptest.NewRequest(
		http.MethodGet, "/recipe/1/comment/1", strings.NewReader(string(reqJson)),
	)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %s", UserJWT2))

	rec := httptest.NewRecorder()

	c := TestE.NewContext(req, rec)
	c.SetPath("/recipe/1/comment/add")
	c.SetParamNames("recipe_id")
	c.SetParamValues("1")

	if assert.NoError(t, TestJwtMiddleware(TestServer.CreateCommentHandle)(c)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)

		respJson := DefaultResponse{}
		err := json.Unmarshal(rec.Body.Bytes(), &respJson)
		assert.Nil(t, err)

		assert.Equal(t, "Оценка пользователя должна быть от 0 до 5", respJson.Message)
	}
}

func TestCreateCommentByAuthor(t *testing.T) {

	reqMap := map[string]interface{}{
		"text": "",
		"rate": 3,
	}
	reqJson, _ := json.Marshal(reqMap)

	req := httptest.NewRequest(
		http.MethodGet, "/recipe/1/comment/1", strings.NewReader(string(reqJson)),
	)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %s", UserJWT))

	rec := httptest.NewRecorder()

	c := TestE.NewContext(req, rec)
	c.SetPath("/recipe/1/comment/add")
	c.SetParamNames("recipe_id")
	c.SetParamValues("1")

	if assert.NoError(t, TestJwtMiddleware(TestServer.CreateCommentHandle)(c)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)

		respJson := DefaultResponse{}
		err := json.Unmarshal(rec.Body.Bytes(), &respJson)
		assert.Nil(t, err)

		assert.Equal(t, "Нельзя оставить отзыв о собственном рецепте", respJson.Message)
	}
}

func TestCreateComment(t *testing.T) {

	reqMap := map[string]interface{}{
		"text": "",
		"rate": 3,
	}
	reqJson, _ := json.Marshal(reqMap)

	req := httptest.NewRequest(
		http.MethodGet, "/recipe/1/comment/1", strings.NewReader(string(reqJson)),
	)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %s", UserJWT2))

	rec := httptest.NewRecorder()

	c := TestE.NewContext(req, rec)
	c.SetPath("/recipe/1/comment/add")
	c.SetParamNames("recipe_id")
	c.SetParamValues("1")

	if assert.NoError(t, TestJwtMiddleware(TestServer.CreateCommentHandle)(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)

		respJson := DefaultResponse{}
		err := json.Unmarshal(rec.Body.Bytes(), &respJson)
		assert.Nil(t, err)

		assert.Equal(t, "Комментарий создан!", respJson.Message)
	}
}

func TestGetComment(t *testing.T) {

	reqMap := map[string]interface{}{
		"visible": false,
	}
	reqJson, _ := json.Marshal(reqMap)

	req := httptest.NewRequest(
		http.MethodGet, "/recipe/1/comment/1", strings.NewReader(string(reqJson)),
	)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %s", UserJWT2))

	rec := httptest.NewRecorder()

	c := TestE.NewContext(req, rec)
	c.SetPath("/recipe/1/comment/1")
	c.SetParamNames("recipe_id", "comment_id")
	c.SetParamValues("1", "1")

	if assert.NoError(t, TestJwtMiddleware(TestServer.GetCommentHandle)(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)

		respJson := models.Comment{}
		err := json.Unmarshal(rec.Body.Bytes(), &respJson)
		assert.Nil(t, err)

		assert.Equal(t, uint(1), respJson.ID)
	}
}

func TestGetAllCommentFromRecipe(t *testing.T) {

	reqMap := map[string]interface{}{
		"visible": false,
	}
	reqJson, _ := json.Marshal(reqMap)

	req := httptest.NewRequest(
		http.MethodGet, "/recipe/1/comments", strings.NewReader(string(reqJson)),
	)
	// req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %s", UserJWT2))

	rec := httptest.NewRecorder()

	c := TestE.NewContext(req, rec)
	c.SetPath("/recipe/1/comments")
	c.SetParamNames("recipe_id")
	c.SetParamValues("1")

	if assert.NoError(t, TestJwtMiddleware(TestServer.GetCommentsHandle)(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)

		respJson := []models.Comment{}
		err := json.Unmarshal(rec.Body.Bytes(), &respJson)
		assert.Nil(t, err)

		assert.Equal(t, 1, len(respJson))
	}
}

func TestDeleteNotExistsComment(t *testing.T) {

	reqMap := map[string]interface{}{
		"visible": false,
	}
	reqJson, _ := json.Marshal(reqMap)

	req := httptest.NewRequest(
		http.MethodGet, "/recipe/1/comment/100/delete", strings.NewReader(string(reqJson)),
	)
	// req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %s", UserJWT2))

	rec := httptest.NewRecorder()

	c := TestE.NewContext(req, rec)
	c.SetPath("/recipe/1/comment/100/delete")
	c.SetParamNames("recipe_id", "comment_id")
	c.SetParamValues("1", "100")

	if assert.NoError(t, TestJwtMiddleware(TestServer.DeleteCommentHandle)(c)) {
		assert.Equal(t, http.StatusNotFound, rec.Code)

		respJson := DefaultResponse{}
		err := json.Unmarshal(rec.Body.Bytes(), &respJson)
		assert.Nil(t, err)

		assert.Equal(t, "Комментарий не найден", respJson.Message)
	}
}

func TestDeleteComment(t *testing.T) {

	reqMap := map[string]interface{}{
		"visible": false,
	}
	reqJson, _ := json.Marshal(reqMap)

	req := httptest.NewRequest(
		http.MethodGet, "/recipe/1/comment/1/delete", strings.NewReader(string(reqJson)),
	)
	// req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %s", UserJWT2))

	rec := httptest.NewRecorder()

	c := TestE.NewContext(req, rec)
	c.SetPath("/recipe/1/comment/1/delete")
	c.SetParamNames("recipe_id", "comment_id")
	c.SetParamValues("1", "1")

	if assert.NoError(t, TestJwtMiddleware(TestServer.DeleteCommentHandle)(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)

		respJson := DefaultResponse{}
		err := json.Unmarshal(rec.Body.Bytes(), &respJson)
		assert.Nil(t, err)

		assert.Equal(t, "Комментарий удален", respJson.Message)
	}
}

func TestCreateIngredient(t *testing.T) {

	reqMap := map[string]interface{}{
		"name":          "a",
		"calories":      0,
		"proteins":      0,
		"fats":          0,
		"carbohydrates": 0,
	}
	reqJson, _ := json.Marshal(reqMap)

	req := httptest.NewRequest(
		http.MethodPost, "/ingredient/create", strings.NewReader(string(reqJson)),
	)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %s", UserJWT))

	rec := httptest.NewRecorder()

	c := TestE.NewContext(req, rec)
	// c.SetPath("/recipe/1/ingredient/add")
	// c.SetParamNames("recipe_id")
	// c.SetParamValues("1")

	if assert.NoError(t, TestJwtMiddleware(TestServer.NewIngredient)(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)

		respJson := DefaultResponse{}
		err := json.Unmarshal(rec.Body.Bytes(), &respJson)
		assert.Nil(t, err)

		assert.Equal(t, "Ингредиент создан!", respJson.Message)
	}
}

func TestCreateStage(t *testing.T) {

	reqMap := map[string]interface{}{
		"description": "a",
	}
	reqJson, _ := json.Marshal(reqMap)

	req := httptest.NewRequest(
		http.MethodPost, "/recipe/1/stage/add", strings.NewReader(string(reqJson)),
	)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %s", UserJWT))

	rec := httptest.NewRecorder()

	c := TestE.NewContext(req, rec)
	c.SetPath("/recipe/1/stage/add")
	c.SetParamNames("recipe_id")
	c.SetParamValues("1")

	if assert.NoError(t, TestJwtMiddleware(TestServer.CreateStageHandle)(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)

		respJson := StageResponse{}
		err := json.Unmarshal(rec.Body.Bytes(), &respJson)
		assert.Nil(t, err)

		assert.Equal(t, "Создан новый этап", respJson.Message)
	}
}

func TestUpdateStage(t *testing.T) {

	reqMap := map[string]interface{}{
		"description": "aa",
	}
	reqJson, _ := json.Marshal(reqMap)

	req := httptest.NewRequest(
		http.MethodPost, "/recipe/1/stage/1/update", strings.NewReader(string(reqJson)),
	)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %s", UserJWT))

	rec := httptest.NewRecorder()

	c := TestE.NewContext(req, rec)
	c.SetPath("/recipe/1/stage/1/update")
	c.SetParamNames("recipe_id", "stage_id")
	c.SetParamValues("1", "1")

	if assert.NoError(t, TestJwtMiddleware(TestServer.UpdateStageHandle)(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)

		respJson := StageResponse{}
		err := json.Unmarshal(rec.Body.Bytes(), &respJson)
		assert.Nil(t, err)

		assert.Equal(t, "Этап обновлен", respJson.Message)
	}
}

func TestDeleteStage(t *testing.T) {

	reqMap := map[string]interface{}{}
	reqJson, _ := json.Marshal(reqMap)

	req := httptest.NewRequest(
		http.MethodPost, "/recipe/1/stage/1/delete", strings.NewReader(string(reqJson)),
	)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %s", UserJWT))

	rec := httptest.NewRecorder()

	c := TestE.NewContext(req, rec)
	c.SetPath("/recipe/1/stage/1/delete")
	c.SetParamNames("recipe_id", "stage_id")
	c.SetParamValues("1", "1")

	if assert.NoError(t, TestJwtMiddleware(TestServer.DeleteStageHandle)(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)

		respJson := StageResponse{}
		err := json.Unmarshal(rec.Body.Bytes(), &respJson)
		assert.Nil(t, err)

		assert.Equal(t, "Этап удален", respJson.Message)
	}
}

func TestDownloadFile(t *testing.T) {

	reqMap := map[string]interface{}{
		"description": "a",
	}
	reqJson, _ := json.Marshal(reqMap)

	req := httptest.NewRequest(
		http.MethodPost, "/assets/a.hehe", strings.NewReader(string(reqJson)),
	)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %s", UserJWT))

	rec := httptest.NewRecorder()

	c := TestE.NewContext(req, rec)
	c.SetPath("/assets/a.hehe")
	c.SetParamNames("filename")
	c.SetParamValues("a.hehe")

	if assert.NoError(t, TestJwtMiddleware(TestServer.DownloadFile)(c)) {
		assert.Equal(t, http.StatusNotFound, rec.Code)

		respJson := DefaultResponse{}
		err := json.Unmarshal(rec.Body.Bytes(), &respJson)
		assert.Nil(t, err)

		assert.Equal(t, "Файл не найден", respJson.Message)
	}
}

func TestDeleteRecipe(t *testing.T) {

	reqMap := map[string]interface{}{
		"name": "b",
	}
	reqJson, _ := json.Marshal(reqMap)

	req := httptest.NewRequest(
		http.MethodGet, "/my-recipe/delete/1", strings.NewReader(string(reqJson)),
	)
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %s", UserJWT))

	rec := httptest.NewRecorder()

	c := TestE.NewContext(req, rec)
	c.SetPath("/my-recipe/delete/1")
	c.SetParamNames("id")
	c.SetParamValues("1")

	if assert.NoError(t, TestJwtMiddleware(TestServer.DeleteRecipeHandle)(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)

		respJson := DefaultResponse{}
		err := json.Unmarshal(rec.Body.Bytes(), &respJson)
		assert.Nil(t, err)
		assert.Equal(t, "Рецепт удален", respJson.Message)

		var recipe models.Recipe
		TestServer.DB.First(&recipe, "id = ?", 1)
		assert.NotEqual(t, "", recipe.DeletedAt)
	}
}

func TestDeleteUser(t *testing.T) {

	reqMap := map[string]interface{}{
		"old_password":     "a",
		"password":         "b",
		"confirm_password": "b",
	}
	reqJson, _ := json.Marshal(reqMap)

	req := httptest.NewRequest(
		http.MethodPost, "/profile/delete", strings.NewReader(string(reqJson)),
	)
	// req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %s", UserJWT2))

	rec := httptest.NewRecorder()

	c := TestE.NewContext(req, rec)

	if assert.NoError(t, TestJwtMiddleware(TestServer.DeleteProfileHandle)(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)

		respJson := DefaultResponse{}
		err := json.Unmarshal(rec.Body.Bytes(), &respJson)
		assert.Nil(t, err)

		assert.Equal(t, "Пользователь успешно удалён!", respJson.Message)

		var user models.User
		TestServer.DB.First(&user, "id = ?", 1)
		assert.NotEqual(t, "", user.DeletedAt)
	}
}

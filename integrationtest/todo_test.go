package integrationtest

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/Hidayathamir/golang-clean-architecture/internal/delivery/http/response"
	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCreateTodo(t *testing.T) {
	TestLogin(t)

	user := GetFirstUser(t)

	requestBody := model.CreateTodoRequest{
		Title:       "Buy groceries",
		Description: "Milk, bread, eggs",
	}
	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	req := httptest.NewRequest(http.MethodPost, "/api/todos", strings.NewReader(string(bodyJson)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", user.Token)

	res, err := app.Test(req)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(res.Body)
	assert.Nil(t, err)

	responseBody := new(response.WebResponse[model.TodoResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, requestBody.Title, responseBody.Data.Title)
	assert.Equal(t, requestBody.Description, responseBody.Data.Description)
	assert.False(t, responseBody.Data.IsCompleted)
	assert.Nil(t, responseBody.Data.CompletedAt)
	assert.NotEmpty(t, responseBody.Data.ID)
	assert.NotZero(t, responseBody.Data.CreatedAt)
	assert.NotZero(t, responseBody.Data.UpdatedAt)
}

func TestCreateTodoFailed(t *testing.T) {
	TestLogin(t)

	user := GetFirstUser(t)

	requestBody := model.CreateTodoRequest{
		Title:       "",
		Description: "",
	}
	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	req := httptest.NewRequest(http.MethodPost, "/api/todos", strings.NewReader(string(bodyJson)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", user.Token)

	res, err := app.Test(req)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(res.Body)
	assert.Nil(t, err)

	responseBody := new(response.WebResponse[model.TodoResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	assert.NotNil(t, responseBody.ErrorMessage)
}

func TestGetTodo(t *testing.T) {
	TestCreateTodo(t)

	user := GetFirstUser(t)
	todo := GetFirstTodo(t, user)

	req := httptest.NewRequest(http.MethodGet, "/api/todos/"+todo.ID, nil)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", user.Token)

	res, err := app.Test(req)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(res.Body)
	assert.Nil(t, err)

	responseBody := new(response.WebResponse[model.TodoResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, todo.ID, responseBody.Data.ID)
	assert.Equal(t, todo.Title, responseBody.Data.Title)
	assert.Equal(t, todo.Description, responseBody.Data.Description)
	assert.Equal(t, todo.IsCompleted, responseBody.Data.IsCompleted)
	assert.Equal(t, todo.CompletedAt, responseBody.Data.CompletedAt)
	assert.Equal(t, todo.CreatedAt, responseBody.Data.CreatedAt)
	assert.Equal(t, todo.UpdatedAt, responseBody.Data.UpdatedAt)
}

func TestGetTodoFailed(t *testing.T) {
	TestCreateTodo(t)

	user := GetFirstUser(t)

	req := httptest.NewRequest(http.MethodGet, "/api/todos/"+uuid.NewString(), nil)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", user.Token)

	res, err := app.Test(req)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(res.Body)
	assert.Nil(t, err)

	responseBody := new(response.WebResponse[model.TodoResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusNotFound, res.StatusCode)
	assert.NotNil(t, responseBody.ErrorMessage)
}

func TestUpdateTodo(t *testing.T) {
	TestCreateTodo(t)

	user := GetFirstUser(t)
	todo := GetFirstTodo(t, user)

	requestBody := model.UpdateTodoRequest{
		Title:       "Buy more groceries",
		Description: "Add vegetables to the list",
	}
	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	req := httptest.NewRequest(http.MethodPut, "/api/todos/"+todo.ID, strings.NewReader(string(bodyJson)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", user.Token)

	res, err := app.Test(req)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(res.Body)
	assert.Nil(t, err)

	responseBody := new(response.WebResponse[model.TodoResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, requestBody.Title, responseBody.Data.Title)
	assert.Equal(t, requestBody.Description, responseBody.Data.Description)
	assert.False(t, responseBody.Data.IsCompleted)
	assert.Nil(t, responseBody.Data.CompletedAt)
	assert.NotZero(t, responseBody.Data.UpdatedAt)
}

func TestUpdateTodoFailed(t *testing.T) {
	TestCreateTodo(t)

	user := GetFirstUser(t)
	todo := GetFirstTodo(t, user)

	requestBody := model.UpdateTodoRequest{
		Title:       "",
		Description: "",
	}
	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	req := httptest.NewRequest(http.MethodPut, "/api/todos/"+todo.ID, strings.NewReader(string(bodyJson)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", user.Token)

	res, err := app.Test(req)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(res.Body)
	assert.Nil(t, err)

	responseBody := new(response.WebResponse[model.TodoResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	assert.NotNil(t, responseBody.ErrorMessage)
}

func TestDeleteTodo(t *testing.T) {
	TestCreateTodo(t)

	user := GetFirstUser(t)
	todo := GetFirstTodo(t, user)

	req := httptest.NewRequest(http.MethodDelete, "/api/todos/"+todo.ID, nil)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", user.Token)

	res, err := app.Test(req)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(res.Body)
	assert.Nil(t, err)

	responseBody := new(response.WebResponse[bool])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.True(t, responseBody.Data)
}

func TestDeleteTodoFailed(t *testing.T) {
	TestCreateTodo(t)

	user := GetFirstUser(t)

	req := httptest.NewRequest(http.MethodDelete, "/api/todos/"+uuid.NewString(), nil)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", user.Token)

	res, err := app.Test(req)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(res.Body)
	assert.Nil(t, err)

	responseBody := new(response.WebResponse[bool])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusNotFound, res.StatusCode)
	assert.NotNil(t, responseBody.ErrorMessage)
}

func TestCompleteTodo(t *testing.T) {
	TestCreateTodo(t)

	user := GetFirstUser(t)
	todo := GetFirstTodo(t, user)

	req := httptest.NewRequest(http.MethodPatch, "/api/todos/"+todo.ID+"/_complete", nil)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", user.Token)

	res, err := app.Test(req)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(res.Body)
	assert.Nil(t, err)

	responseBody := new(response.WebResponse[model.TodoResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.True(t, responseBody.Data.IsCompleted)
	assert.NotNil(t, responseBody.Data.CompletedAt)

	var updated entity.Todo
	err = db.Where("id = ?", todo.ID).First(&updated).Error
	assert.Nil(t, err)
	assert.True(t, updated.IsCompleted)
	assert.NotNil(t, updated.CompletedAt)
}

func TestListTodos(t *testing.T) {
	TestLogin(t)

	user := GetFirstUser(t)

	CreateTodos(t, user, 15)

	req := httptest.NewRequest(http.MethodGet, "/api/todos", nil)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", user.Token)

	res, err := app.Test(req)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(res.Body)
	assert.Nil(t, err)

	responseBody := new(response.WebResponse[model.TodoResponseList])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, 10, len(responseBody.Data))
	assert.Equal(t, int64(15), responseBody.Paging.TotalItem)
	assert.Equal(t, int64(2), responseBody.Paging.TotalPage)
	assert.Equal(t, 1, responseBody.Paging.Page)
	assert.Equal(t, 10, responseBody.Paging.Size)
}

func TestListTodosWithPagination(t *testing.T) {
	TestLogin(t)

	user := GetFirstUser(t)

	CreateTodos(t, user, 25)

	req := httptest.NewRequest(http.MethodGet, "/api/todos?page=2&size=10", nil)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", user.Token)

	res, err := app.Test(req)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(res.Body)
	assert.Nil(t, err)

	responseBody := new(response.WebResponse[model.TodoResponseList])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, 10, len(responseBody.Data))
	assert.Equal(t, int64(25), responseBody.Paging.TotalItem)
	assert.Equal(t, int64(3), responseBody.Paging.TotalPage)
	assert.Equal(t, 2, responseBody.Paging.Page)
	assert.Equal(t, 10, responseBody.Paging.Size)
}

func TestListTodosWithFilters(t *testing.T) {
	TestLogin(t)

	user := GetFirstUser(t)

	CreateTodos(t, user, 5)

	now := time.Now().UnixMilli()
	for i := range 3 {
		completedAt := now + int64(i)
		todo := &entity.Todo{
			ID:          uuid.NewString(),
			UserID:      user.ID,
			Title:       "Todo Completed " + strconv.Itoa(i),
			Description: "Completed task " + strconv.Itoa(i),
			IsCompleted: true,
			CompletedAt: &completedAt,
		}
		err := db.Create(todo).Error
		assert.Nil(t, err)
	}

	req := httptest.NewRequest(http.MethodGet, "/api/todos?title=Todo&is_completed=true", nil)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", user.Token)

	res, err := app.Test(req)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(res.Body)
	assert.Nil(t, err)

	responseBody := new(response.WebResponse[model.TodoResponseList])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, 3, len(responseBody.Data))
	assert.Equal(t, int64(3), responseBody.Paging.TotalItem)
	assert.Equal(t, int64(1), responseBody.Paging.TotalPage)
	assert.Equal(t, 1, responseBody.Paging.Page)
	assert.Equal(t, 10, responseBody.Paging.Size)

	for _, todo := range responseBody.Data {
		assert.True(t, todo.IsCompleted)
		assert.NotNil(t, todo.CompletedAt)
	}
}

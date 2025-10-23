package http

import (
	"math"
	"net/http"
	"strconv"

	"github.com/Hidayathamir/golang-clean-architecture/internal/delivery/http/middleware"
	"github.com/Hidayathamir/golang-clean-architecture/internal/delivery/http/response"
	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
	todousecase "github.com/Hidayathamir/golang-clean-architecture/internal/usecase/todo"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type TodoController struct {
	Usecase todousecase.TodoUsecase
	Log     *logrus.Logger
}

func NewTodoController(usecase todousecase.TodoUsecase, log *logrus.Logger) *TodoController {
	return &TodoController{
		Usecase: usecase,
		Log:     log,
	}
}

// Create godoc
//
//	@Summary		Create todo
//	@Description	Create a new todo
//	@Tags			todos
//	@Security		SimpleApiKeyAuth
//	@Param			request	body		model.CreateTodoRequest	true	"Create Todo Request"
//	@Success		200		{object}	response.WebResponse[model.TodoResponse]
//	@Router			/api/todos [post]
func (c *TodoController) Create(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)

	req := new(model.CreateTodoRequest)
	if err := ctx.BodyParser(req); err != nil {
		err = errkit.BadRequest(err)
		return errkit.AddFuncName("http.(*TodoController).Create", err)
	}

	req.UserID = auth.ID

	res, err := c.Usecase.Create(ctx.UserContext(), req)
	if err != nil {
		return errkit.AddFuncName("http.(*TodoController).Create", err)
	}

	return response.Data(ctx, http.StatusOK, res)
}

// List godoc
//
//	@Summary		List todos
//	@Description	List todos with optional filters
//	@Tags			todos
//	@Security		SimpleApiKeyAuth
//	@Param			title			query		string	false	"Filter by title"
//	@Param			is_completed	query		bool	false	"Filter by completion state"
//	@Param			page			query		int		false	"Page number"	default(1)
//	@Param			size			query		int		false	"Page size"		default(10)
//	@Success		200				{object}	response.WebResponse[[]model.TodoResponse]
//	@Router			/api/todos [get]
func (c *TodoController) List(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)

	var isCompleted *bool
	if raw := ctx.Query("is_completed"); raw != "" {
		value, err := strconv.ParseBool(raw)
		if err != nil {
			err = errkit.BadRequest(err)
			return errkit.AddFuncName("http.(*TodoController).List", err)
		}
		isCompleted = &value
	}

	req := &model.ListTodoRequest{
		UserID:      auth.ID,
		Title:       ctx.Query("title", ""),
		IsCompleted: isCompleted,
		Page:        ctx.QueryInt("page", 1),
		Size:        ctx.QueryInt("size", 10),
	}

	res, total, err := c.Usecase.List(ctx.UserContext(), req)
	if err != nil {
		return errkit.AddFuncName("http.(*TodoController).List", err)
	}

	paging := &response.PageMetadata{
		Page:      req.Page,
		Size:      req.Size,
		TotalItem: total,
		TotalPage: int64(math.Ceil(float64(total) / float64(req.Size))),
	}

	return response.DataPaging(ctx, http.StatusOK, res, paging)
}

// Get godoc
//
//	@Summary		Get todo
//	@Description	Get todo by ID
//	@Tags			todos
//	@Security		SimpleApiKeyAuth
//	@Param			todoId	path		string	true	"Todo ID"
//	@Success		200		{object}	response.WebResponse[model.TodoResponse]
//	@Router			/api/todos/{todoId} [get]
func (c *TodoController) Get(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)

	req := &model.GetTodoRequest{
		UserID: auth.ID,
		ID:     ctx.Params("todoId"),
	}

	res, err := c.Usecase.Get(ctx.UserContext(), req)
	if err != nil {
		return errkit.AddFuncName("http.(*TodoController).Get", err)
	}

	return response.Data(ctx, http.StatusOK, res)
}

// Update godoc
//
//	@Summary		Update todo
//	@Description	Update todo attributes
//	@Tags			todos
//	@Security		SimpleApiKeyAuth
//	@Param			todoId	path		string					true	"Todo ID"
//	@Param			request	body		model.UpdateTodoRequest	true	"Update Todo Request"
//	@Success		200		{object}	response.WebResponse[model.TodoResponse]
//	@Router			/api/todos/{todoId} [put]
func (c *TodoController) Update(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)

	req := new(model.UpdateTodoRequest)
	if err := ctx.BodyParser(req); err != nil {
		err = errkit.BadRequest(err)
		return errkit.AddFuncName("http.(*TodoController).Update", err)
	}

	req.UserID = auth.ID
	req.ID = ctx.Params("todoId")

	res, err := c.Usecase.Update(ctx.UserContext(), req)
	if err != nil {
		return errkit.AddFuncName("http.(*TodoController).Update", err)
	}

	return response.Data(ctx, http.StatusOK, res)
}

// Delete godoc
//
//	@Summary		Delete todo
//	@Description	Delete todo by ID
//	@Tags			todos
//	@Security		SimpleApiKeyAuth
//	@Param			todoId	path		string	true	"Todo ID"
//	@Success		200		{object}	response.WebResponse[bool]
//	@Router			/api/todos/{todoId} [delete]
func (c *TodoController) Delete(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)

	req := &model.DeleteTodoRequest{
		UserID: auth.ID,
		ID:     ctx.Params("todoId"),
	}

	if err := c.Usecase.Delete(ctx.UserContext(), req); err != nil {
		return errkit.AddFuncName("http.(*TodoController).Delete", err)
	}

	return response.Data(ctx, http.StatusOK, true)
}

// Complete godoc
//
//	@Summary		Complete todo
//	@Description	Mark todo as completed and emit event
//	@Tags			todos
//	@Security		SimpleApiKeyAuth
//	@Param			todoId	path		string	true	"Todo ID"
//	@Success		200		{object}	response.WebResponse[model.TodoResponse]
//	@Router			/api/todos/{todoId}/_complete [patch]
func (c *TodoController) Complete(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)

	req := &model.CompleteTodoRequest{
		UserID: auth.ID,
		ID:     ctx.Params("todoId"),
	}

	res, err := c.Usecase.Complete(ctx.UserContext(), req)
	if err != nil {
		return errkit.AddFuncName("http.(*TodoController).Complete", err)
	}

	return response.Data(ctx, http.StatusOK, res)
}

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
	"github.com/Hidayathamir/golang-clean-architecture/pkg/telemetry"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

type TodoController struct {
	Config  *viper.Viper
	Usecase todousecase.TodoUsecase
}

func NewTodoController(cfg *viper.Viper, usecase todousecase.TodoUsecase) *TodoController {
	return &TodoController{
		Config:  cfg,
		Usecase: usecase,
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
	span := telemetry.StartController(ctx)
	defer span.End()

	auth := middleware.GetUser(ctx)

	req := new(model.CreateTodoRequest)
	if err := ctx.BodyParser(req); err != nil {
		err = errkit.BadRequest(err)
		return errkit.AddFuncName(err)
	}

	req.UserID = auth.ID

	res, err := c.Usecase.Create(ctx.UserContext(), req)
	if err != nil {
		return errkit.AddFuncName(err)
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
//	@Success		200				{object}	response.WebResponse[model.TodoResponseList]
//	@Router			/api/todos [get]
func (c *TodoController) List(ctx *fiber.Ctx) error {
	span := telemetry.StartController(ctx)
	defer span.End()

	auth := middleware.GetUser(ctx)

	var isCompleted *bool
	if raw := ctx.Query("is_completed"); raw != "" {
		value, err := strconv.ParseBool(raw)
		if err != nil {
			err = errkit.BadRequest(err)
			return errkit.AddFuncName(err)
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
		return errkit.AddFuncName(err)
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
//	@Param			todoId	path		int	true	"Todo ID"
//	@Success		200		{object}	response.WebResponse[model.TodoResponse]
//	@Router			/api/todos/{todoId} [get]
func (c *TodoController) Get(ctx *fiber.Ctx) error {
	span := telemetry.StartController(ctx)
	defer span.End()

	auth := middleware.GetUser(ctx)

	todoID, err := strconv.ParseInt(ctx.Params("todoId"), 10, 64)
	if err != nil {
		err = errkit.BadRequest(err)
		return errkit.AddFuncName(err)
	}

	req := &model.GetTodoRequest{
		UserID: auth.ID,
		ID:     todoID,
	}

	res, err := c.Usecase.Get(ctx.UserContext(), req)
	if err != nil {
		return errkit.AddFuncName(err)
	}

	return response.Data(ctx, http.StatusOK, res)
}

// Update godoc
//
//	@Summary		Update todo
//	@Description	Update todo attributes
//	@Tags			todos
//	@Security		SimpleApiKeyAuth
//	@Param			todoId	path		int					true	"Todo ID"
//	@Param			request	body		model.UpdateTodoRequest	true	"Update Todo Request"
//	@Success		200		{object}	response.WebResponse[model.TodoResponse]
//	@Router			/api/todos/{todoId} [put]
func (c *TodoController) Update(ctx *fiber.Ctx) error {
	span := telemetry.StartController(ctx)
	defer span.End()

	auth := middleware.GetUser(ctx)

	req := new(model.UpdateTodoRequest)
	if err := ctx.BodyParser(req); err != nil {
		err = errkit.BadRequest(err)
		return errkit.AddFuncName(err)
	}

	todoID, err := strconv.ParseInt(ctx.Params("todoId"), 10, 64)
	if err != nil {
		err = errkit.BadRequest(err)
		return errkit.AddFuncName(err)
	}

	req.UserID = auth.ID
	req.ID = todoID

	res, err := c.Usecase.Update(ctx.UserContext(), req)
	if err != nil {
		return errkit.AddFuncName(err)
	}

	return response.Data(ctx, http.StatusOK, res)
}

// Delete godoc
//
//	@Summary		Delete todo
//	@Description	Delete todo by ID
//	@Tags			todos
//	@Security		SimpleApiKeyAuth
//	@Param			todoId	path		int	true	"Todo ID"
//	@Success		200		{object}	response.WebResponse[bool]
//	@Router			/api/todos/{todoId} [delete]
func (c *TodoController) Delete(ctx *fiber.Ctx) error {
	span := telemetry.StartController(ctx)
	defer span.End()

	auth := middleware.GetUser(ctx)

	todoID, err := strconv.ParseInt(ctx.Params("todoId"), 10, 64)
	if err != nil {
		err = errkit.BadRequest(err)
		return errkit.AddFuncName(err)
	}

	req := &model.DeleteTodoRequest{
		UserID: auth.ID,
		ID:     todoID,
	}

	if err := c.Usecase.Delete(ctx.UserContext(), req); err != nil {
		return errkit.AddFuncName(err)
	}

	return response.Data(ctx, http.StatusOK, true)
}

// Complete godoc
//
//	@Summary		Complete todo
//	@Description	Mark todo as completed and emit event
//	@Tags			todos
//	@Security		SimpleApiKeyAuth
//	@Param			todoId	path		int	true	"Todo ID"
//	@Success		200		{object}	response.WebResponse[model.TodoResponse]
//	@Router			/api/todos/{todoId}/_complete [patch]
func (c *TodoController) Complete(ctx *fiber.Ctx) error {
	span := telemetry.StartController(ctx)
	defer span.End()

	auth := middleware.GetUser(ctx)

	todoID, err := strconv.ParseInt(ctx.Params("todoId"), 10, 64)
	if err != nil {
		err = errkit.BadRequest(err)
		return errkit.AddFuncName(err)
	}

	req := &model.CompleteTodoRequest{
		UserID: auth.ID,
		ID:     todoID,
	}

	res, err := c.Usecase.Complete(ctx.UserContext(), req)
	if err != nil {
		return errkit.AddFuncName(err)
	}

	return response.Data(ctx, http.StatusOK, res)
}

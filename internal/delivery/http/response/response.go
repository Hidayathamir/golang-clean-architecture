package response

import (
	"errors"
	"golang-clean-architecture/pkg/httperror"

	"github.com/gofiber/fiber/v2"
)

type WebResponse[T any] struct {
	Data         T             `json:"data"`
	Paging       *PageMetadata `json:"paging,omitempty"`
	ErrorID      string        `json:"error_id,omitempty"`
	ErrorMessage string        `json:"error_message,omitempty"`
}

type PageResponse[T any] struct {
	Data         []T          `json:"data,omitempty"`
	PageMetadata PageMetadata `json:"paging,omitempty"`
}

type PageMetadata struct {
	Page      int   `json:"page"`
	Size      int   `json:"size"`
	TotalItem int64 `json:"total_item"`
	TotalPage int64 `json:"total_page"`
}

func Data(ctx *fiber.Ctx, status int, data any) error {
	res := WebResponse[any]{}
	res.Data = data
	return ctx.Status(status).JSON(res)
}

func DataPaging(ctx *fiber.Ctx, status int, data any, paging *PageMetadata) error {
	res := WebResponse[any]{}
	res.Data = data
	res.Paging = paging
	return ctx.Status(status).JSON(res)
}

func Error(ctx *fiber.Ctx, err error) error {
	if err == nil {
		return nil
	}

	httpErr := LoadErrAsHTTPError(err)

	res := WebResponse[any]{}
	res.ErrorID = httpErr.ID
	res.ErrorMessage = httpErr.Message

	return ctx.Status(httpErr.HTTPCode).JSON(res)
}

// LoadErrAsHTTPError load error as HTTPError with use InternalServerError as default.
func LoadErrAsHTTPError(err error) *httperror.HTTPError {
	httpErr := httperror.InternalServerError()
	errors.As(err, &httpErr)

	return httpErr
}

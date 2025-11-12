package response

import (
	"net/http"

	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/gofiber/fiber/v2"
)

type WebResponse[T any] struct {
	Data         T             `json:"data"`
	Paging       *PageMetadata `json:"paging,omitempty"`
	ErrorMessage string        `json:"error_message"`
	ErrorDetail  []string      `json:"error_detail"`
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
		return ctx.Status(http.StatusInternalServerError).SendString("something wrong")
	}

	httpErr := errkit.GetHTTPError(err)

	res := WebResponse[any]{}
	res.ErrorMessage = httpErr.Message
	res.ErrorDetail = errkit.Split(err)

	return ctx.Status(httpErr.HTTPCode).JSON(res)
}

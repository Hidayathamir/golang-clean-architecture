package http

import (
	"net/http"

	"github.com/Hidayathamir/golang-clean-architecture/internal/delivery/http/response"
	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
	"github.com/Hidayathamir/golang-clean-architecture/internal/usecase/image"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/telemetry"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

type ImageController struct {
	Config  *viper.Viper
	Usecase image.ImageUsecase
}

func NewImageController(cfg *viper.Viper, useCase image.ImageUsecase) *ImageController {
	return &ImageController{
		Config:  cfg,
		Usecase: useCase,
	}
}

// Upload godoc
//
//	@Summary		Upload image
//	@Description	Upload an image file
//	@Tags			images
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			image	formData	file	true	"Image to upload"
//	@Security		SimpleApiKeyAuth
//	@Success		200	{object}	response.WebResponse[model.ImageResponse]
//	@Router			/api/images [post]
func (c *ImageController) Upload(ctx *fiber.Ctx) error {
	span := telemetry.StartController(ctx)
	defer span.End()

	file, err := ctx.FormFile("image")
	if err != nil {
		err = errkit.BadRequest(err)
		return errkit.AddFuncName(err)
	}

	req := new(model.UploadImageRequest)
	req.File = file

	res, err := c.Usecase.Upload(ctx.UserContext(), req)
	if err != nil {
		return errkit.AddFuncName(err)
	}

	return response.Data(ctx, http.StatusOK, res)
}

// Like godoc
//
//	@Summary		Like image
//	@Description	Like an image
//	@Tags			images
//	@Accept			json
//	@Produce		json
//	@Param			request	body	model.LikeImageRequest	true	"Like Image Request"
//	@Security		SimpleApiKeyAuth
//	@Success		200	{object}	response.WebResponse[string]
//	@Router			/api/images/_like [post]
func (c *ImageController) Like(ctx *fiber.Ctx) error {
	span := telemetry.StartController(ctx)
	defer span.End()

	req := new(model.LikeImageRequest)
	err := ctx.BodyParser(req)
	if err != nil {
		err = errkit.BadRequest(err)
		return errkit.AddFuncName(err)
	}

	err = c.Usecase.Like(ctx.UserContext(), req)
	if err != nil {
		return errkit.AddFuncName(err)
	}

	return response.Data(ctx, http.StatusOK, "ok")
}

// Comment godoc
//
//	@Summary		Comment image
//	@Description	Comment an image
//	@Tags			images
//	@Accept			json
//	@Produce		json
//	@Param			request	body	model.CommentImageRequest	true	"Comment Image Request"
//	@Security		SimpleApiKeyAuth
//	@Success		200	{object}	response.WebResponse[string]
//	@Router			/api/images/_comment [post]
func (c *ImageController) Comment(ctx *fiber.Ctx) error {
	span := telemetry.StartController(ctx)
	defer span.End()

	req := new(model.CommentImageRequest)
	err := ctx.BodyParser(req)
	if err != nil {
		err = errkit.BadRequest(err)
		return errkit.AddFuncName(err)
	}

	err = c.Usecase.Comment(ctx.UserContext(), req)
	if err != nil {
		return errkit.AddFuncName(err)
	}

	return response.Data(ctx, http.StatusOK, "ok")
}

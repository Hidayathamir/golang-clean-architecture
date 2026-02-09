package http

import (
	"net/http"
	"strconv"

	"github.com/Hidayathamir/golang-clean-architecture/internal/config"
	"github.com/Hidayathamir/golang-clean-architecture/internal/delivery/http/response"
	"github.com/Hidayathamir/golang-clean-architecture/internal/dto"
	"github.com/Hidayathamir/golang-clean-architecture/internal/usecase/image"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/telemetry"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/x"
	"github.com/gofiber/fiber/v2"
)

type ImageController struct {
	Cfg     *config.Config
	Usecase image.ImageUsecase
}

func NewImageController(cfg *config.Config, usecase image.ImageUsecase) *ImageController {
	return &ImageController{
		Cfg:     cfg,
		Usecase: usecase,
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
//	@Success		200	{object}	response.WebResponse[dto.ImageResponse]
//	@Router			/api/images [post]
func (c *ImageController) Upload(ctx *fiber.Ctx) error {
	span := telemetry.StartController(ctx)
	defer span.End()

	file, err := ctx.FormFile("image")
	if err != nil {
		err = errkit.BadRequest(err)
		x.Logger.WithContext(ctx.UserContext()).WithError(err).Error()
		return errkit.AddFuncName(err)
	}

	req := new(dto.UploadImageRequest)
	req.File = file

	res, err := c.Usecase.Upload(ctx.UserContext(), req)
	if err != nil {
		x.Logger.WithContext(ctx.UserContext()).WithError(err).Error()
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
//	@Param			request	body	dto.LikeImageRequest	true	"Like Image Request"
//	@Security		SimpleApiKeyAuth
//	@Success		200	{object}	response.WebResponse[string]
//	@Router			/api/images/_like [post]
func (c *ImageController) Like(ctx *fiber.Ctx) error {
	span := telemetry.StartController(ctx)
	defer span.End()

	req := new(dto.LikeImageRequest)
	err := ctx.BodyParser(req)
	if err != nil {
		err = errkit.BadRequest(err)
		x.Logger.WithContext(ctx.UserContext()).WithError(err).Error()
		return errkit.AddFuncName(err)
	}

	err = c.Usecase.Like(ctx.UserContext(), req)
	if err != nil {
		x.Logger.WithContext(ctx.UserContext()).WithError(err).Error()
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
//	@Param			request	body	dto.CommentImageRequest	true	"Comment Image Request"
//	@Security		SimpleApiKeyAuth
//	@Success		200	{object}	response.WebResponse[string]
//	@Router			/api/images/_comment [post]
func (c *ImageController) Comment(ctx *fiber.Ctx) error {
	span := telemetry.StartController(ctx)
	defer span.End()

	req := new(dto.CommentImageRequest)
	err := ctx.BodyParser(req)
	if err != nil {
		err = errkit.BadRequest(err)
		x.Logger.WithContext(ctx.UserContext()).WithError(err).Error()
		return errkit.AddFuncName(err)
	}

	err = c.Usecase.Comment(ctx.UserContext(), req)
	if err != nil {
		x.Logger.WithContext(ctx.UserContext()).WithError(err).Error()
		return errkit.AddFuncName(err)
	}

	return response.Data(ctx, http.StatusOK, "ok")
}

// GetLike godoc
//
//	@Summary		Get image likes
//	@Description	Get list of likes for an image
//	@Tags			images
//	@Produce		json
//	@Param			imageId	path	int	true	"Image ID"
//	@Security		SimpleApiKeyAuth
//	@Success		200	{object}	response.WebResponse[dto.LikeResponseList]
//	@Router			/api/images/{imageId}/likes [get]
func (c *ImageController) GetLike(ctx *fiber.Ctx) error {
	span := telemetry.StartController(ctx)
	defer span.End()

	imageID, err := strconv.ParseInt(ctx.Params("imageId"), 10, 64)
	if err != nil {
		err = errkit.BadRequest(err)
		x.Logger.WithContext(ctx.UserContext()).WithError(err).Error()
		return errkit.AddFuncName(err)
	}

	req := &dto.GetLikeRequest{
		ImageID: imageID,
	}

	res, err := c.Usecase.GetLike(ctx.UserContext(), req)
	if err != nil {
		x.Logger.WithContext(ctx.UserContext()).WithError(err).Error()
		return errkit.AddFuncName(err)
	}

	return response.Data(ctx, http.StatusOK, res)
}

// GetComment godoc
//
//	@Summary		Get image comments
//	@Description	Get list of comments for an image
//	@Tags			images
//	@Produce		json
//	@Param			imageId	path	int	true	"Image ID"
//	@Security		SimpleApiKeyAuth
//	@Success		200	{object}	response.WebResponse[dto.CommentResponseList]
//	@Router			/api/images/{imageId}/comments [get]
func (c *ImageController) GetComment(ctx *fiber.Ctx) error {
	span := telemetry.StartController(ctx)
	defer span.End()

	imageID, err := strconv.ParseInt(ctx.Params("imageId"), 10, 64)
	if err != nil {
		err = errkit.BadRequest(err)
		x.Logger.WithContext(ctx.UserContext()).WithError(err).Error()
		return errkit.AddFuncName(err)
	}

	req := &dto.GetCommentRequest{
		ImageID: imageID,
	}

	res, err := c.Usecase.GetComment(ctx.UserContext(), req)
	if err != nil {
		x.Logger.WithContext(ctx.UserContext()).WithError(err).Error()
		return errkit.AddFuncName(err)
	}

	return response.Data(ctx, http.StatusOK, res)
}

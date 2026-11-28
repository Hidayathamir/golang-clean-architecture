package http

import (
	"net/http"
	"strconv"

	"github.com/Hidayathamir/golang-clean-architecture/internal/config"
	"github.com/Hidayathamir/golang-clean-architecture/internal/dto"
	"github.com/Hidayathamir/golang-clean-architecture/internal/inbound/http/response"
	"github.com/Hidayathamir/golang-clean-architecture/internal/usecase/imageusecase"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/logkit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/telemetry"
	"github.com/gofiber/fiber/v2"
)

type ImageController struct {
	Cfg     *config.Config
	Usecase imageusecase.ImageUsecase
}

func NewImageController(cfg *config.Config, usecase imageusecase.ImageUsecase) *ImageController {
	return &ImageController{
		Cfg:     cfg,
		Usecase: usecase,
	}
}

// PresignUpload godoc
//
//	@Summary		Presign image upload
//	@Description	Returns a presigned S3 PUT URL. The client uploads the image bytes directly to S3, then calls /api/images/_confirm with the key.
//	@Tags			images
//	@Accept			json
//	@Produce		json
//	@Param			request	body	dto.PresignUploadImageRequest	true	"Presign Upload Request"
//	@Security		SimpleApiKeyAuth
//	@Success		200	{object}	response.WebResponse[dto.PresignUploadImageResponse]
//	@Router			/api/images [post]
func (c *ImageController) PresignUpload(ctx *fiber.Ctx) error {
	span := telemetry.StartController(ctx)
	defer span.End()

	req := dto.PresignUploadImageRequest{}
	err := ctx.BodyParser(&req)
	if err != nil {
		err = errkit.SetCode(err, http.StatusBadRequest)
		logkit.Logger.WithContext(ctx.UserContext()).WithError(err).Error()
		return errkit.AddFuncName(err, "http.(*ImageController).PresignUpload")
	}

	res, err := c.Usecase.PresignUpload(ctx.UserContext(), req)
	if err != nil {
		logkit.Logger.WithContext(ctx.UserContext()).WithError(err).Error()
		return errkit.AddFuncName(err, "http.(*ImageController).PresignUpload")
	}

	return response.Data(ctx, http.StatusOK, res)
}

// ConfirmUpload godoc
//
//	@Summary		Confirm image upload
//	@Description	Finalizes a client-side upload after the bytes landed in S3. Creates the image record and emits the upload event.
//	@Tags			images
//	@Accept			json
//	@Produce		json
//	@Param			request	body	dto.ConfirmUploadImageRequest	true	"Confirm Upload Request"
//	@Security		SimpleApiKeyAuth
//	@Success		200	{object}	response.WebResponse[dto.ImageResponse]
//	@Router			/api/images/_confirm [post]
func (c *ImageController) ConfirmUpload(ctx *fiber.Ctx) error {
	span := telemetry.StartController(ctx)
	defer span.End()

	req := dto.ConfirmUploadImageRequest{}
	err := ctx.BodyParser(&req)
	if err != nil {
		err = errkit.SetCode(err, http.StatusBadRequest)
		logkit.Logger.WithContext(ctx.UserContext()).WithError(err).Error()
		return errkit.AddFuncName(err, "http.(*ImageController).ConfirmUpload")
	}

	res, err := c.Usecase.ConfirmUpload(ctx.UserContext(), req)
	if err != nil {
		logkit.Logger.WithContext(ctx.UserContext()).WithError(err).Error()
		return errkit.AddFuncName(err, "http.(*ImageController).ConfirmUpload")
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

	req := dto.LikeImageRequest{}
	err := ctx.BodyParser(&req)
	if err != nil {
		err = errkit.SetCode(err, http.StatusBadRequest)
		logkit.Logger.WithContext(ctx.UserContext()).WithError(err).Error()
		return errkit.AddFuncName(err, "http.(*ImageController).Like")
	}

	err = c.Usecase.Like(ctx.UserContext(), req)
	if err != nil {
		logkit.Logger.WithContext(ctx.UserContext()).WithError(err).Error()
		return errkit.AddFuncName(err, "http.(*ImageController).Like")
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

	req := dto.CommentImageRequest{}
	err := ctx.BodyParser(&req)
	if err != nil {
		err = errkit.SetCode(err, http.StatusBadRequest)
		logkit.Logger.WithContext(ctx.UserContext()).WithError(err).Error()
		return errkit.AddFuncName(err, "http.(*ImageController).Comment")
	}

	err = c.Usecase.Comment(ctx.UserContext(), req)
	if err != nil {
		logkit.Logger.WithContext(ctx.UserContext()).WithError(err).Error()
		return errkit.AddFuncName(err, "http.(*ImageController).Comment")
	}

	return response.Data(ctx, http.StatusOK, "ok")
}

// GetImage godoc
//
//	@Summary		Get image
//	@Description	Get a single image by id
//	@Tags			images
//	@Produce		json
//	@Param			imageId	path	int	true	"Image ID"
//	@Security		SimpleApiKeyAuth
//	@Success		200	{object}	response.WebResponse[dto.ImageResponse]
//	@Router			/api/images/{imageId} [get]
func (c *ImageController) GetImage(ctx *fiber.Ctx) error {
	span := telemetry.StartController(ctx)
	defer span.End()

	imageID, err := strconv.ParseInt(ctx.Params("imageId"), 10, 64)
	if err != nil {
		err = errkit.SetCode(err, http.StatusBadRequest)
		logkit.Logger.WithContext(ctx.UserContext()).WithError(err).Error()
		return errkit.AddFuncName(err, "http.(*ImageController).GetImage")
	}

	req := dto.GetImageRequest{
		ID: imageID,
	}

	res, err := c.Usecase.GetImage(ctx.UserContext(), req)
	if err != nil {
		logkit.Logger.WithContext(ctx.UserContext()).WithError(err).Error()
		return errkit.AddFuncName(err, "http.(*ImageController).GetImage")
	}

	return response.Data(ctx, http.StatusOK, res)
}

// GetFeed godoc
//
//	@Summary		Get feed
//	@Description	Get posts from users the current user follows, newest first (cursor paginated)
//	@Tags			images
//	@Produce		json
//	@Param			limit	query	int		false	"Page size (default 20, max 100)"
//	@Param			cursor	query	string	false	"Opaque cursor from the previous page; empty for the first page"
//	@Security		SimpleApiKeyAuth
//	@Success		200	{object}	response.WebResponse[dto.FeedResponse]
//	@Router			/api/feed [get]
func (c *ImageController) GetFeed(ctx *fiber.Ctx) error {
	span := telemetry.StartController(ctx)
	defer span.End()

	req := dto.FeedRequest{}
	if limitStr := ctx.Query("limit"); limitStr != "" {
		limit, err := strconv.Atoi(limitStr)
		if err != nil {
			err = errkit.SetCode(err, http.StatusBadRequest)
			logkit.Logger.WithContext(ctx.UserContext()).WithError(err).Error()
			return errkit.AddFuncName(err, "http.(*ImageController).GetFeed")
		}
		req.Page.Limit = limit
	}
	req.Page.Cursor = ctx.Query("cursor")

	res, err := c.Usecase.GetFeed(ctx.UserContext(), req)
	if err != nil {
		logkit.Logger.WithContext(ctx.UserContext()).WithError(err).Error()
		return errkit.AddFuncName(err, "http.(*ImageController).GetFeed")
	}

	return response.Data(ctx, http.StatusOK, res)
}

// GetUserImages godoc
//
//	@Summary		Get user images
//	@Description	Get a user's post grid, newest first (cursor paginated)
//	@Tags			images
//	@Produce		json
//	@Param			id		path	int		true	"User ID"
//	@Param			limit	query	int		false	"Page size (default 20, max 100)"
//	@Param			cursor	query	string	false	"Opaque cursor from the previous page; empty for the first page"
//	@Security		SimpleApiKeyAuth
//	@Success		200	{object}	response.WebResponse[dto.UserImagesResponse]
//	@Router			/api/users/{id}/images [get]
func (c *ImageController) GetUserImages(ctx *fiber.Ctx) error {
	span := telemetry.StartController(ctx)
	defer span.End()

	userID, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
	if err != nil {
		err = errkit.SetCode(err, http.StatusBadRequest)
		logkit.Logger.WithContext(ctx.UserContext()).WithError(err).Error()
		return errkit.AddFuncName(err, "http.(*ImageController).GetUserImages")
	}

	req := dto.UserImagesRequest{
		UserID: userID,
	}
	if limitStr := ctx.Query("limit"); limitStr != "" {
		limit, err := strconv.Atoi(limitStr)
		if err != nil {
			err = errkit.SetCode(err, http.StatusBadRequest)
			logkit.Logger.WithContext(ctx.UserContext()).WithError(err).Error()
			return errkit.AddFuncName(err, "http.(*ImageController).GetUserImages")
		}
		req.Page.Limit = limit
	}
	req.Page.Cursor = ctx.Query("cursor")

	res, err := c.Usecase.GetUserImages(ctx.UserContext(), req)
	if err != nil {
		logkit.Logger.WithContext(ctx.UserContext()).WithError(err).Error()
		return errkit.AddFuncName(err, "http.(*ImageController).GetUserImages")
	}

	return response.Data(ctx, http.StatusOK, res)
}

// Explore godoc
//
//	@Summary		Explore images
//	@Description	Discover posts via Elasticsearch, newest first (cursor paginated)
//	@Tags			images
//	@Produce		json
//	@Param			limit	query	int		false	"Page size (default 20, max 100)"
//	@Param			cursor	query	string	false	"Opaque cursor from the previous page; empty for the first page"
//	@Security		SimpleApiKeyAuth
//	@Success		200	{object}	response.WebResponse[dto.ExploreResponse]
//	@Router			/api/explore [get]
func (c *ImageController) Explore(ctx *fiber.Ctx) error {
	span := telemetry.StartController(ctx)
	defer span.End()

	req := dto.ExploreRequest{}
	if limitStr := ctx.Query("limit"); limitStr != "" {
		limit, err := strconv.Atoi(limitStr)
		if err != nil {
			err = errkit.SetCode(err, http.StatusBadRequest)
			logkit.Logger.WithContext(ctx.UserContext()).WithError(err).Error()
			return errkit.AddFuncName(err, "http.(*ImageController).Explore")
		}
		req.Page.Limit = limit
	}
	req.Page.Cursor = ctx.Query("cursor")

	res, err := c.Usecase.Explore(ctx.UserContext(), req)
	if err != nil {
		logkit.Logger.WithContext(ctx.UserContext()).WithError(err).Error()
		return errkit.AddFuncName(err, "http.(*ImageController).Explore")
	}

	return response.Data(ctx, http.StatusOK, res)
}

// DeleteImage godoc
//
//	@Summary		Delete image
//	@Description	Delete an image (only the owner may delete)
//	@Tags			images
//	@Produce		json
//	@Param			imageId	path	int	true	"Image ID"
//	@Security		SimpleApiKeyAuth
//	@Success		200	{object}	response.WebResponse[string]
//	@Router			/api/images/{imageId} [delete]
func (c *ImageController) DeleteImage(ctx *fiber.Ctx) error {
	span := telemetry.StartController(ctx)
	defer span.End()

	imageID, err := strconv.ParseInt(ctx.Params("imageId"), 10, 64)
	if err != nil {
		err = errkit.SetCode(err, http.StatusBadRequest)
		logkit.Logger.WithContext(ctx.UserContext()).WithError(err).Error()
		return errkit.AddFuncName(err, "http.(*ImageController).DeleteImage")
	}

	req := dto.DeleteImageRequest{
		ID: imageID,
	}

	err = c.Usecase.DeleteImage(ctx.UserContext(), req)
	if err != nil {
		logkit.Logger.WithContext(ctx.UserContext()).WithError(err).Error()
		return errkit.AddFuncName(err, "http.(*ImageController).DeleteImage")
	}

	return response.Data(ctx, http.StatusOK, "ok")
}

// UpdateImage godoc
//
//	@Summary		Update image
//	@Description	Edit an image caption (only the owner may edit)
//	@Tags			images
//	@Accept			json
//	@Produce		json
//	@Param			imageId	path	int						true	"Image ID"
//	@Param			request	body	dto.UpdateImageRequest	true	"Update Image Request"
//	@Security		SimpleApiKeyAuth
//	@Success		200	{object}	response.WebResponse[string]
//	@Router			/api/images/{imageId} [patch]
func (c *ImageController) UpdateImage(ctx *fiber.Ctx) error {
	span := telemetry.StartController(ctx)
	defer span.End()

	imageID, err := strconv.ParseInt(ctx.Params("imageId"), 10, 64)
	if err != nil {
		err = errkit.SetCode(err, http.StatusBadRequest)
		logkit.Logger.WithContext(ctx.UserContext()).WithError(err).Error()
		return errkit.AddFuncName(err, "http.(*ImageController).UpdateImage")
	}

	req := dto.UpdateImageRequest{}
	err = ctx.BodyParser(&req)
	if err != nil {
		err = errkit.SetCode(err, http.StatusBadRequest)
		logkit.Logger.WithContext(ctx.UserContext()).WithError(err).Error()
		return errkit.AddFuncName(err, "http.(*ImageController).UpdateImage")
	}
	req.ID = imageID

	err = c.Usecase.UpdateImage(ctx.UserContext(), req)
	if err != nil {
		logkit.Logger.WithContext(ctx.UserContext()).WithError(err).Error()
		return errkit.AddFuncName(err, "http.(*ImageController).UpdateImage")
	}

	return response.Data(ctx, http.StatusOK, "ok")
}

// Unlike godoc
//
//	@Summary		Unlike image
//	@Description	Remove the current user's like from an image
//	@Tags			images
//	@Produce		json
//	@Param			imageId	path	int	true	"Image ID"
//	@Security		SimpleApiKeyAuth
//	@Success		200	{object}	response.WebResponse[string]
//	@Router			/api/images/{imageId}/like [delete]
func (c *ImageController) Unlike(ctx *fiber.Ctx) error {
	span := telemetry.StartController(ctx)
	defer span.End()

	imageID, err := strconv.ParseInt(ctx.Params("imageId"), 10, 64)
	if err != nil {
		err = errkit.SetCode(err, http.StatusBadRequest)
		logkit.Logger.WithContext(ctx.UserContext()).WithError(err).Error()
		return errkit.AddFuncName(err, "http.(*ImageController).Unlike")
	}

	req := dto.UnlikeImageRequest{
		ImageID: imageID,
	}

	err = c.Usecase.Unlike(ctx.UserContext(), req)
	if err != nil {
		logkit.Logger.WithContext(ctx.UserContext()).WithError(err).Error()
		return errkit.AddFuncName(err, "http.(*ImageController).Unlike")
	}

	return response.Data(ctx, http.StatusOK, "ok")
}

// DeleteImageComment godoc
//
//	@Summary		Delete image comment
//	@Description	Delete a comment (comment owner or image owner)
//	@Tags			images
//	@Produce		json
//	@Param			imageId		path	int	true	"Image ID"
//	@Param			commentId	path	int	true	"Comment ID"
//	@Security		SimpleApiKeyAuth
//	@Success		200	{object}	response.WebResponse[string]
//	@Router			/api/images/{imageId}/comments/{commentId} [delete]
func (c *ImageController) DeleteImageComment(ctx *fiber.Ctx) error {
	span := telemetry.StartController(ctx)
	defer span.End()

	imageID, err := strconv.ParseInt(ctx.Params("imageId"), 10, 64)
	if err != nil {
		err = errkit.SetCode(err, http.StatusBadRequest)
		logkit.Logger.WithContext(ctx.UserContext()).WithError(err).Error()
		return errkit.AddFuncName(err, "http.(*ImageController).DeleteImageComment")
	}

	commentID, err := strconv.ParseInt(ctx.Params("commentId"), 10, 64)
	if err != nil {
		err = errkit.SetCode(err, http.StatusBadRequest)
		logkit.Logger.WithContext(ctx.UserContext()).WithError(err).Error()
		return errkit.AddFuncName(err, "http.(*ImageController).DeleteImageComment")
	}

	req := dto.DeleteImageCommentRequest{
		ImageID:   imageID,
		CommentID: commentID,
	}

	err = c.Usecase.DeleteImageComment(ctx.UserContext(), req)
	if err != nil {
		logkit.Logger.WithContext(ctx.UserContext()).WithError(err).Error()
		return errkit.AddFuncName(err, "http.(*ImageController).DeleteImageComment")
	}

	return response.Data(ctx, http.StatusOK, "ok")
}

// GetLike godoc
//
//	@Summary		Get image likes
//	@Description	Get list of likes for an image (cursor paginated)
//	@Tags			images
//	@Produce		json
//	@Param			imageId	path	int		true	"Image ID"
//	@Param			limit	query	int		false	"Page size (default 20, max 100)"
//	@Param			cursor	query	string	false	"Opaque cursor from the previous page; empty for the first page"
//	@Security		SimpleApiKeyAuth
//	@Success		200	{object}	response.WebResponse[dto.LikesResponse]
//	@Router			/api/images/{imageId}/likes [get]
func (c *ImageController) GetLike(ctx *fiber.Ctx) error {
	span := telemetry.StartController(ctx)
	defer span.End()

	imageID, err := strconv.ParseInt(ctx.Params("imageId"), 10, 64)
	if err != nil {
		err = errkit.SetCode(err, http.StatusBadRequest)
		logkit.Logger.WithContext(ctx.UserContext()).WithError(err).Error()
		return errkit.AddFuncName(err, "http.(*ImageController).GetLike")
	}

	req := dto.GetLikeRequest{
		ImageID: imageID,
	}
	if limitStr := ctx.Query("limit"); limitStr != "" {
		limit, err := strconv.Atoi(limitStr)
		if err != nil {
			err = errkit.SetCode(err, http.StatusBadRequest)
			logkit.Logger.WithContext(ctx.UserContext()).WithError(err).Error()
			return errkit.AddFuncName(err, "http.(*ImageController).GetLike")
		}
		req.Page.Limit = limit
	}
	req.Page.Cursor = ctx.Query("cursor")

	res, err := c.Usecase.GetLike(ctx.UserContext(), req)
	if err != nil {
		logkit.Logger.WithContext(ctx.UserContext()).WithError(err).Error()
		return errkit.AddFuncName(err, "http.(*ImageController).GetLike")
	}

	return response.Data(ctx, http.StatusOK, res)
}

// GetComment godoc
//
//	@Summary		Get image comments
//	@Description	Get list of comments for an image (cursor paginated)
//	@Tags			images
//	@Produce		json
//	@Param			imageId	path	int		true	"Image ID"
//	@Param			limit	query	int		false	"Page size (default 20, max 100)"
//	@Param			cursor	query	string	false	"Opaque cursor from the previous page; empty for the first page"
//	@Security		SimpleApiKeyAuth
//	@Success		200	{object}	response.WebResponse[dto.CommentsResponse]
//	@Router			/api/images/{imageId}/comments [get]
func (c *ImageController) GetComment(ctx *fiber.Ctx) error {
	span := telemetry.StartController(ctx)
	defer span.End()

	imageID, err := strconv.ParseInt(ctx.Params("imageId"), 10, 64)
	if err != nil {
		err = errkit.SetCode(err, http.StatusBadRequest)
		logkit.Logger.WithContext(ctx.UserContext()).WithError(err).Error()
		return errkit.AddFuncName(err, "http.(*ImageController).GetComment")
	}

	req := dto.GetCommentRequest{
		ImageID: imageID,
	}
	if limitStr := ctx.Query("limit"); limitStr != "" {
		limit, err := strconv.Atoi(limitStr)
		if err != nil {
			err = errkit.SetCode(err, http.StatusBadRequest)
			logkit.Logger.WithContext(ctx.UserContext()).WithError(err).Error()
			return errkit.AddFuncName(err, "http.(*ImageController).GetComment")
		}
		req.Page.Limit = limit
	}
	req.Page.Cursor = ctx.Query("cursor")

	res, err := c.Usecase.GetComment(ctx.UserContext(), req)
	if err != nil {
		logkit.Logger.WithContext(ctx.UserContext()).WithError(err).Error()
		return errkit.AddFuncName(err, "http.(*ImageController).GetComment")
	}

	return response.Data(ctx, http.StatusOK, res)
}

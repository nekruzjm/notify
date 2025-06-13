package event

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"notifications/internal/api/resp"
	"notifications/internal/api/resp/code"
	"notifications/internal/service/admin"
	"notifications/pkg/util/strset"
)

// UploadImage
//
//	@Summary		Upload image
//	@Description	Provider event id - `{id}` and multi-lang image key - `{language}` as api route var to upload image for event
//	@Tags			Events
//	@Accept			multipart/form-data
//	@Produce		application/json
//	@Param			image	formData	file								true	"Load images with png, jpg, jpeg extensions"
//	@Success		200		{object}	resp.Response{payload=eventModel}	"Success"
//	@Failure		400		{object}	resp.Response						"Bad request"
//	@Failure		401		{object}	resp.Response						"Invalid authorization data"
//	@Failure		404		{object}	resp.Response						"Not found"
//	@Failure		500		{object}	resp.Response						"Internal Error"
//	@Router			/notifications-internal/v1/events/{id}/image/{language} [post]
func (h *handler) UploadImage(c *gin.Context) {
	var (
		ctx      = c.Request.Context()
		id       = strset.ToInt(c.Param(_id))
		lang     = c.Param(_language)
		response resp.Response
	)

	defer resp.JSON(c.Writer, code.Success, &response)

	adminUser, ok := ctx.Value(admin.CtxKey).(admin.Admin)
	if !ok {
		response = resp.RespondErr(resp.ErrUnauthorized)
		return
	}

	file, fileHeader, err := c.Request.FormFile(_image)
	if err != nil {
		if errors.Is(err, http.ErrMissingFile) {
			err = resp.Wrap(resp.ErrBadRequest, "file is required")
		}
		response = resp.RespondErr(err)
		return
	}
	defer func() { _ = file.Close() }()

	serviceResponse, err := h.service.UploadImage(ctx, adminUser, lang, id, file, fileHeader)
	if err != nil {
		response = resp.RespondErr(err)
		return
	}

	response = resp.Success
	response.Payload = serviceResponse
}

// RemoveImage
//
//	@Summary		Remove image
//	@Description	Provider event id - `{id}` and multi-lang image key - `{language}` as api route var to remove specific image of event
//	@Tags			Events
//	@Produce		application/json
//	@Success		200	{object}	resp.Response{payload=eventModel}	"Success"
//	@Failure		400	{object}	resp.Response						"Bad request"
//	@Failure		401	{object}	resp.Response						"Invalid authorization data"
//	@Failure		404	{object}	resp.Response						"Not found"
//	@Failure		500	{object}	resp.Response						"Internal Error"
//	@Router			/notifications-internal/v1/events/{id}/image/{language} [delete]
func (h *handler) RemoveImage(c *gin.Context) {
	var (
		ctx      = c.Request.Context()
		id       = strset.ToInt(c.Param(_id))
		lang     = c.Param(_language)
		response resp.Response
	)

	defer resp.JSON(c.Writer, code.Success, &response)

	adminUser, ok := ctx.Value(admin.CtxKey).(admin.Admin)
	if !ok {
		response = resp.RespondErr(resp.ErrUnauthorized)
		return
	}

	serviceResponse, err := h.service.RemoveImage(ctx, adminUser, lang, id)
	if err != nil {
		response = resp.RespondErr(err)
		return
	}

	response = resp.Success
	response.Payload = serviceResponse
}

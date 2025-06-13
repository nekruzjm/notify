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

// LoadUsers
//
//	@Summary	Subscribe list of users to event
//	@Tags		Events
//	@Accept		multipart/form-data
//	@Produce	application/json
//	@Param		users	formData	file								true	"CSV file with `userID` header and data. Do not provide other headers"
//	@Success	202		{object}	resp.Response{payload=eventModel}	"Accepted"
//	@Failure	400		{object}	resp.Response						"Bad request"
//	@Failure	401		{object}	resp.Response						"Invalid authorization data"
//	@Failure	404		{object}	resp.Response						"Not found"
//	@Failure	500		{object}	resp.Response						"Internal Error"
//	@Router		/notifications-internal/v1/events/{id}/load-users [post]
func (h *handler) LoadUsers(c *gin.Context) {
	var (
		ctx      = c.Request.Context()
		id       = strset.ToInt(c.Param(_id))
		response resp.Response
	)

	defer resp.JSON(c.Writer, code.Success, &response)

	adminUser, ok := ctx.Value(admin.CtxKey).(admin.Admin)
	if !ok {
		response = resp.RespondErr(resp.ErrUnauthorized)
		return
	}

	file, fileHeader, err := c.Request.FormFile(_users)
	if err != nil {
		if errors.Is(err, http.ErrMissingFile) {
			err = resp.Wrap(resp.ErrBadRequest, "file is required")
		}
		response = resp.RespondErr(err)
		return
	}
	defer func() { _ = file.Close() }()

	serviceResponse, err := h.service.LoadUsers(ctx, adminUser, id, file, fileHeader)
	if err != nil {
		response = resp.RespondErr(err)
		return
	}

	response = resp.Accepted
	response.Payload = serviceResponse
}

// LoadAllUsers
//
//	@Summary	Subscribe all users to the event
//	@Tags		Events
//	@Produce	application/json
//	@Success	202	{object}	resp.Response{payload=eventModel}	"Accepted"
//	@Failure	400	{object}	resp.Response						"Bad request"
//	@Failure	401	{object}	resp.Response						"Invalid authorization data"
//	@Failure	404	{object}	resp.Response						"Not found"
//	@Failure	500	{object}	resp.Response						"Internal Error"
//	@Router		/notifications-internal/v1/events/{id}/load-all-users [post]
func (h *handler) LoadAllUsers(c *gin.Context) {
	var (
		ctx      = c.Request.Context()
		id       = strset.ToInt(c.Param(_id))
		response resp.Response
	)

	defer resp.JSON(c.Writer, code.Success, &response)

	adminUser, ok := ctx.Value(admin.CtxKey).(admin.Admin)
	if !ok {
		response = resp.RespondErr(resp.ErrUnauthorized)
		return
	}

	serviceResponse, err := h.service.LoadAllUsers(ctx, adminUser, id)
	if err != nil {
		response = resp.RespondErr(err)
		return
	}

	response = resp.Accepted
	response.Payload = serviceResponse
}

package event

import (
	"github.com/gin-gonic/gin"

	"notifications/internal/api/resp"
	"notifications/internal/api/resp/code"
	"notifications/internal/service/admin"
	"notifications/pkg/util/strset"
)

// Run
//
//	@Summary	Run event manually
//	@Tags		Events
//	@Produce	application/json
//	@Success	200	{object}	resp.Response{payload=eventModel}	"Success"
//	@Failure	400	{object}	resp.Response						"Bad request"
//	@Failure	401	{object}	resp.Response						"Invalid authorization data"
//	@Failure	404	{object}	resp.Response						"Not found"
//	@Failure	500	{object}	resp.Response						"Internal Error"
//	@Router		/notifications-internal/v1/events/{id}/run [post]
func (h *handler) Run(c *gin.Context) {
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

	serviceResponse, err := h.service.RunEvent(ctx, adminUser, id)
	if err != nil {
		response = resp.RespondErr(err)
		return
	}

	response = resp.Success
	response.Payload = serviceResponse
}

package push

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"notifications/internal/api/resp"
	"notifications/internal/api/resp/code"
	"notifications/internal/service/push"
	"notifications/pkg/util/serializer"
)

// Send
// @Description	All fields except `personExternalRef` (crm_client_id) are required.
// @Description	- If you want to send push with `personExternalRef`, do not provide `phone`.
// @Description	- If `showInFeed` is true, the push will be shown in the feed; otherwise, it will be hidden.
// @Description	- If the users status is inactive or their push setting is disabled, the push will be saved in the feed but not sent to the device.
// @Description	In that case, the payload will be `inactive_user#fake_message_id` or `disabled_push#fake_message_id`.
// @Tags			External
// @Accept			application/json
// @Produce		application/json
// @Param			X-UserId		header		string			true	"Provide user ID created on the server side"
// @Param			X-RequestId		header		string			true	"Provide unique request ID to build hash and track the request"
// @Param			X-Date			header		string			true	"Provide current date in RFC1123 format (e.g., Mon, 02 Jan 2006 15:04:05 MST) to build hash"
// @Param			X-UserAction	header		string			true	"Provide user action (push, sms) to send push"
// @Param			X-RequestDigest	header		string			true	"Provide hash sum built with HMAC-SHA256 from the `X-Date:X-RequestId` using the secret key created on the server side"
// @Param			data			body		externalRequest	true	"Request payload"
// @Success		200				{object}	resp.Response	"Success"
// @Failure		400				{object}	resp.Response	"Bad request"
// @Failure		401				{object}	resp.Response	"Invalid authorization data"
// @Failure		404				{object}	resp.Response	"Not found"
// @Failure		500				{object}	resp.Response	"Internal Error"
// @Security		SignatureAuth
// @Router			/notifications-external/v1/push [post]
func (h *handler) Send(c *gin.Context) {
	var (
		ctx          = c.Request.Context()
		requestID, _ = ctx.Value(_requestID).(string)
		apiClient, _ = ctx.Value(_apiClient).(string)
		response     resp.Response
		request      externalRequest
	)

	defer resp.JSON(c.Writer, code.Success, &response)

	err := serializer.BodyToJSON(c.Request, &request)
	if err != nil {
		err = resp.Wrap(resp.ErrBadRequest, err.Error())
		response = resp.RespondErr(err)
		return
	}

	h.logger.Info("send external push",
		zap.Any("request", request),
		zap.String(_apiClient, apiClient),
		zap.String(_requestID, requestID))

	var message = new(push.Request)
	message.ExternalRequest.ID = requestID
	message.ExternalRequest.APIClient = apiClient
	message.ExternalRequest.PersonExternalRef = request.PersonExternalRef
	message.ExternalRequest.Phone = request.Phone
	message.ExternalRequest.Title = request.Title
	message.ExternalRequest.Body = request.Body
	message.ExternalRequest.PushType = request.PushType
	message.ShowInFeed = request.ShowInFeed
	message.IsInternal = false

	messageID, err := h.service.Send(ctx, message)
	if err != nil {
		response = resp.RespondErr(err)
		return
	}

	h.logger.Info("sent external push",
		zap.String(_apiClient, apiClient),
		zap.String(_requestID, requestID),
		zap.String("serviceResponse", messageID))

	response = resp.Success
	response.Payload = messageID
}

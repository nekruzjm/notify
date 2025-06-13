package event

import (
	"github.com/gin-gonic/gin"

	"notifications/internal/api/resp"
	"notifications/internal/api/resp/code"
	"notifications/internal/service/admin"
	"notifications/internal/service/event"
	"notifications/pkg/util/serializer"
	"notifications/pkg/util/strset"
)

// Get
//
//	@Summary	Get list of events
//	@Tags		Events
//	@Accept		application/json
//	@Produce	application/json
//	@Param		id		query		string								false	"apply filter with id"
//	@Param		status	query		string								false	"apply filter with status"
//	@Param		topic	query		string								false	"apply filter with topic"
//	@Param		limit	query		string								false	"apply filter with limit, 10 settled by default"
//	@Param		offset	query		string								false	"apply filter with offset, 0 settled by default"
//	@Success	200		{object}	resp.Response{payload=[]eventModel}	"Success"
//	@Failure	401		{object}	resp.Response						"Invalid authorization data"
//	@Failure	404		{object}	resp.Response						"List not found"
//	@Failure	500		{object}	resp.Response						"Internal Error"
//	@Router		/notifications-internal/v1/events [get]
func (h *handler) Get(c *gin.Context) {
	var (
		ctx      = c.Request.Context()
		id       = strset.ToInt(c.Query(_id))
		status   = c.Query(_status)
		topic    = c.Query(_topic)
		limit    = strset.ToInt(c.Query(_limit))
		offset   = strset.ToInt(c.Query(_offset))
		response resp.Response
	)

	defer resp.JSON(c.Writer, code.Success, &response)

	serviceResponse, err := h.service.GetEvents(ctx, event.Filter{
		ID:     uint(id),
		Topic:  topic,
		Status: status,
		Limit:  uint(limit),
		Offset: uint(offset),
	})
	if err != nil {
		response = resp.RespondErr(err)
		return
	}

	response = resp.Success
	response.Payload = serviceResponse
}

// Create
//
//	@Summary	Create event
//	@Tags		Events
//	@Accept		application/json
//	@Produce	application/json
//	@Param		data	body		request								true	"Request"
//	@Success	200		{object}	resp.Response{payload=eventModel}	"Success"
//	@Failure	400		{object}	resp.Response						"Bad request"
//	@Failure	401		{object}	resp.Response						"Invalid authorization data"
//	@Failure	500		{object}	resp.Response						"Internal Error"
//	@Router		/notifications-internal/v1/events [post]
func (h *handler) Create(c *gin.Context) {
	var (
		ctx      = c.Request.Context()
		response resp.Response
		r        request
	)

	defer resp.JSON(c.Writer, code.Success, &response)

	adminUser, ok := ctx.Value(admin.CtxKey).(admin.Admin)
	if !ok {
		response = resp.RespondErr(resp.ErrUnauthorized)
		return
	}

	err := serializer.BodyToJSON(c.Request, &r)
	if err != nil {
		err = resp.Wrap(resp.ErrBadRequest, err.Error())
		response = resp.RespondErr(err)
		return
	}

	var eventRequest = &event.Request{
		Status:      r.Status,
		Topic:       r.Topic,
		Category:    r.Category,
		Link:        r.Link,
		ScheduledAt: r.ScheduledAt,
		ExtraData:   r.ExtraData,
		Title:       r.Title,
		Body:        r.Body,
	}

	serviceResponse, err := h.service.Create(ctx, adminUser, eventRequest)
	if err != nil {
		response = resp.RespondErr(err)
		return
	}

	response = resp.Success
	response.Payload = serviceResponse
}

// Update
//
//	@Summary	Update event
//	@Tags		Events
//	@Accept		application/json
//	@Produce	application/json
//	@Param		data	body		request								true	"Request"
//	@Success	200		{object}	resp.Response{payload=eventModel}	"Success"
//	@Failure	400		{object}	resp.Response						"Bad request"
//	@Failure	401		{object}	resp.Response						"Invalid authorization data"
//	@Failure	404		{object}	resp.Response						"Not found"
//	@Failure	500		{object}	resp.Response						"Internal Error"
//	@Router		/notifications-internal/v1/events/{id} [put]
func (h *handler) Update(c *gin.Context) {
	var (
		ctx      = c.Request.Context()
		id       = strset.ToInt(c.Param(_id))
		response resp.Response
		r        request
	)

	defer resp.JSON(c.Writer, code.Success, &response)

	adminUser, ok := ctx.Value(admin.CtxKey).(admin.Admin)
	if !ok {
		response = resp.RespondErr(resp.ErrUnauthorized)
		return
	}

	err := serializer.BodyToJSON(c.Request, &r)
	if err != nil {
		err = resp.Wrap(resp.ErrBadRequest, err.Error())
		response = resp.RespondErr(err)
		return
	}

	var eventRequest = &event.Request{
		ID:          id,
		Status:      r.Status,
		Topic:       r.Topic,
		Category:    r.Category,
		Link:        r.Link,
		ScheduledAt: r.ScheduledAt,
		ExtraData:   r.ExtraData,
		Title:       r.Title,
		Body:        r.Body,
	}

	serviceResponse, err := h.service.Update(ctx, adminUser, eventRequest)
	if err != nil {
		response = resp.RespondErr(err)
		return
	}

	response = resp.Success
	response.Payload = serviceResponse
}

// Delete
//
//	@Summary	Delete event
//	@Tags		Events
//	@Accept		application/json
//	@Produce	application/json
//	@Success	200	{object}	resp.Response	"Success"
//	@Failure	401	{object}	resp.Response	"Invalid authorization data"
//	@Failure	404	{object}	resp.Response	"Not found"
//	@Failure	500	{object}	resp.Response	"Internal Error"
//	@Router		/notifications-internal/v1/events/{id} [delete]
func (h *handler) Delete(c *gin.Context) {
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

	err := h.service.Delete(ctx, adminUser, id)
	if err != nil {
		response = resp.RespondErr(err)
		return
	}

	response = resp.Success
}

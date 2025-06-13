package middleware

import (
	"context"
	"errors"
	"slices"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"notifications/internal/api/resp"
	"notifications/internal/api/resp/code"
	"notifications/internal/repo/repomodel"
	"notifications/pkg/lib/security/hasher"
	"notifications/pkg/util/strset"
)

const (
	_digestKey     = "X-RequestDigest"
	_userKey       = "X-UserId"
	_userActionKey = "X-UserAction"
	_dateKey       = "X-Date"
	_requestKey    = "X-RequestId"
	_authorization = "Authorization"
)

func (m *mw) ProtectExternal() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			userID     = c.GetHeader(_userKey)
			userAction = c.GetHeader(_userActionKey)
			requestID  = c.GetHeader(_requestKey)
			date       = c.GetHeader(_dateKey)
			digest     = c.GetHeader(_digestKey)
			clientIP   = c.ClientIP()
			ctx        = context.Background()
			response   = resp.Response{}
		)

		if strset.IsSliceEmpty(userID, requestID, userAction, digest, date) {
			m.logger.Warning("empty headers now allowed",
				zap.String("userID", userID),
				zap.String("userAction", userAction),
				zap.String("clientIP", clientIP),
				zap.String("requestID", requestID))

			response = resp.Unauthorized
			response.Message = "Empty headers are not allowed"
			resp.GinJSONAbort(c, code.Unauthorized, response)
			return
		}

		client, err := m.apiClientRepo.GetByUserID(ctx, userID)
		if err != nil {
			if !errors.Is(err, repomodel.ErrNotFound) {
				m.logger.Error("err occurred during getting api client",
					zap.Error(err),
					zap.String("userID", userID),
					zap.String("userAction", userAction),
					zap.String("clientIP", clientIP),
					zap.String("requestID", requestID))
			}
			response = resp.Unauthorized
			response.Message = "Invalid X-UserId, api client not found"
			resp.GinJSONAbort(c, code.Unauthorized, response)
			return
		}

		// Each client can make 3 requests per second. After a second, the counter will reset
		//limiter := m.rateLimiter.NewSlidingWindowLimiter(userID, 3, time.Second)
		//if !limiter.IsAllowed() {
		//	m.logger.Warning("rate limit exceeded",
		//		zap.String("userID", userID),
		//		zap.String("userAction", userAction),
		//		zap.String("clientIP", clientIP),
		//		zap.String("requestID", requestID))
		//
		//	response = resp.TooManyRequests
		//	response.Message = "Rate limit exceeded, try again later"
		//	resp.GinJSONAbort(c, code.TooManyRequests, response)
		//	return
		//}

		if !slices.Contains(client.Permissions, userAction) {
			m.logger.Warning("user permission denied",
				zap.String("userID", userID),
				zap.String("userAction", userAction),
				zap.String("clientIP", clientIP),
				zap.String("requestID", requestID))

			response = resp.Forbidden
			response.Message = "User doesn't have permission to perform this action"
			resp.GinJSONAbort(c, code.Forbidden, response)
			return
		}

		if _, err = time.Parse(time.RFC1123, date); err != nil {
			response = resp.BadRequest
			response.Message = "Invalid X-Date format, must be RFC1123"
			resp.GinJSONAbort(c, code.BadRequest, response)
			return
		}

		var (
			reqIDCacheKey = ":mw-req-id:" + requestID
			val           int
		)

		if err = m.cache.Get(ctx, reqIDCacheKey, &val); err == nil {
			response = resp.BadRequest
			response.Message = "Invalid X-RequestId, request already exists"
			resp.GinJSONAbort(c, code.BadRequest, response)
			return
		}

		err = m.cache.Set(ctx, reqIDCacheKey, 0, 12*time.Hour)
		if err != nil {
			m.logger.Error("err occurred during setting request id cache", zap.Error(err))
			resp.GinJSONAbort(c, code.InternalErr, resp.InternalErr)
			return
		}

		hash, _ := hasher.GenerateSHA2(client.APIKey, date+":"+requestID)

		if digest != hash {
			m.logger.Warning("digest mismatch",
				zap.String("userID", userID),
				zap.String("userAction", userAction),
				zap.String("clientIP", clientIP),
				zap.String("requestID", requestID))

			response = resp.Unauthorized
			response.Message = "Invalid X-RequestDigest, hash doesn't match"
			resp.GinJSONAbort(c, code.Unauthorized, response)
			return
		}

		c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), "requestID", requestID))
		c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), "apiClient", userID))
		c.Next()
	}
}

func (m *mw) ProtectInternal() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			token    = c.GetHeader(_authorization)
			ctx      = context.Background()
			response resp.Response
		)

		if strset.IsSliceEmpty(token) {
			response = resp.Unauthorized
			response.Message = "Empty headers are not allowed"
			resp.GinJSONAbort(c, code.Unauthorized, response)
			return
		}

		ctx = context.WithValue(ctx, "method", c.Request.Method)
		ctx = context.WithValue(ctx, "url", c.Request.URL.Path)

		admin, err := m.admin.Authorize(ctx, token)
		if err != nil {
			switch {
			case errors.Is(err, resp.ErrForbidden):
				response = resp.Forbidden
			case errors.Is(err, resp.ErrUnauthorized):
				response = resp.Unauthorized
			default:
				response = resp.InternalErr
			}
			resp.GinJSONAbort(c, response.Code, response)
			return
		}

		admin.IP = c.ClientIP()
		c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), "admin", admin))
		c.Next()
	}
}

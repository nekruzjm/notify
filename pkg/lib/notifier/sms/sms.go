package sms

import (
	"context"

	"github.com/imroc/req/v3"
	"go.uber.org/zap"
)

func (s *sms) Send(_ context.Context, request Request) error {
	const _route = "/api/v1/Sms"

	var (
		url    = s.config.GetString("sms.url") + _route
		header = s.headers()
	)

	s.logger.Debug("sending sms", zap.Any("request", request))

	resp, err := req.R().
		SetBody(request).
		SetHeaders(header).
		Post(url)
	if err != nil {
		s.logger.Error("err sending sms", zap.Error(err), zap.String("url", url))
		return err
	}

	s.logger.Debug("sms sent", zap.String("phone", request.Phone), zap.Any("response", resp))

	if resp.IsErrorState() {
		s.logger.Error("incorrect status", zap.String("resp", resp.String()), zap.Error(resp.Err))
		return err
	}

	return nil
}

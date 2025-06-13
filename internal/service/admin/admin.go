package admin

import (
	"context"
	"errors"

	"notifications/internal/api/resp"
	"notifications/internal/gateway/gatewaymodel"
)

func (s *service) Authorize(ctx context.Context, token string) (Admin, error) {
	admin, err := s.gateway.Authorize(ctx, token)
	if err != nil {
		switch {
		case errors.Is(err, gatewaymodel.ErrUnauthorized):
			return Admin{}, resp.ErrUnauthorized
		case errors.Is(err, gatewaymodel.ErrForbidden):
			return Admin{}, resp.ErrForbidden
		default:
			return Admin{}, resp.ErrInternalErr
		}
	}

	return Admin{
		ID:        admin.ID,
		Username:  admin.Username,
		FullName:  admin.FullName,
		CountryID: admin.CountryID,
	}, nil
}

package admin

import (
	"context"
)

func (g *gateway) Authorize(ctx context.Context, token string) (admin Admin, err error) {
	var (
	//reqUrl, _    = ctx.Value("url").(string)
	//reqMethod, _ = ctx.Value("method").(string)
	//url          = g.config.GetString("admin.url") + "/adminusers/validate/token"
	//headers      = map[string]string{
	//	_authorization: token,
	//	_serviceName:   _notifications,
	//	_reqUrl:        reqUrl,
	//	_method:        reqMethod,
	//}
	//response gatewayResponse
	)

	// TODO mock for tests
	return Admin{
		ID:        1,
		Username:  "admin",
		FullName:  "Admin",
		CountryID: 1,
	}, nil

	//resp, err := req.C().R().
	//	SetHeaders(headers).
	//	SetSuccessResult(&response).
	//	Get(url)
	//if err != nil {
	//	g.logger.Error("Error on getting services", zap.Error(err), zap.Any("url", url))
	//	return Admin{}, err
	//}
	//
	//if resp.IsErrorState() {
	//	g.logger.Error("incorrect response status", zap.Error(err), zap.String("response", resp.String()))
	//	return Admin{}, err
	//}
	//
	//const _customUnauthorizedCode = 1510
	//
	//switch response.Code {
	//case http.StatusOK:
	//	return response.Payload, nil
	//case http.StatusUnauthorized, _customUnauthorizedCode:
	//	return Admin{}, gatewaymodel.ErrUnauthorized
	//case http.StatusForbidden:
	//	return Admin{}, gatewaymodel.ErrForbidden
	//default:
	//	return Admin{}, gatewaymodel.ErrInternal
	//}
}

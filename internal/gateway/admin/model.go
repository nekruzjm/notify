package admin

//type gatewayResponse struct {
//	Code    int   `json:"code"`
//	Payload Admin `json:"payload"`
//}

type Admin struct {
	ID        int    `json:"id"`
	CountryID int    `json:"countryID"`
	Username  string `json:"username"`
	FullName  string `json:"fullname"`
}

// auth headers
//const (
//	_authorization = "Authorization"
//	_serviceName   = "Service-Name"
//	_reqUrl        = "ReqUrl"
//	_method        = "Method"
//)

// service name
//const _notifications = "notifications"

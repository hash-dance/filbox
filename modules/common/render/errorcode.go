package render

// ErrorCode define error struct
type ErrorCode struct {
	Message string
	ErrNumber
}

var (
	// Unauthorized 认证失败
	Unauthorized = ErrorCode{"Unauthorized", ErrUnauthorized}
	// NotFound 数据不存在
	NotFound = ErrorCode{"Data Not Found", ErrNotFound}
	// InvalidData 数据格式错误
	InvalidData = ErrorCode{"InvalidData", ErrInvalidData}
	// ServerError 服务器内部错误
	ServerError = ErrorCode{"Server Error", ErrServerError}
	// OauthAuthError oauth 认证错误
	OauthAuthError = ErrorCode{"Oauth Auth Error", ErrOauthAuthError}
	// DenyAccessError 拒绝访问
	DenyAccessError = ErrorCode{"Deny Access", ErrDenyAccessError}
	// RedirectError 重定向
	RedirectError = ErrorCode{"Redirect", ErrRedirect}
	// RouterError 路由错误
	RouterError = ErrorCode{"Router Error", ErrRouterError}
)

// ErrNumber defined errNumber type
type ErrNumber int

var (
	// Success request success
	Success ErrNumber = 0 // 请求成功
	// ErrNotFound request not found error
	ErrNotFound ErrNumber = 1 // 数据不存在
	// ErrInvalidData invalid data form request
	ErrInvalidData ErrNumber = 2 // 数据错误
	// ErrUnauthorized unauthorized
	ErrUnauthorized ErrNumber = 3 // 认证失败
	// ErrServerError server error code
	ErrServerError ErrNumber = 4 // 服务器内部错误
	// ErrOauthAuthError oauth auth error
	ErrOauthAuthError ErrNumber = 5
	// ErrDenyAccessError 拒绝访问
	ErrDenyAccessError ErrNumber = 6
	// ErrRedirect 重定向
	ErrRedirect ErrNumber = 7
	// ErrRouterError 路由错误
	ErrRouterError ErrNumber = 8
)

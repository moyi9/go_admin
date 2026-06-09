// Package apperror 定义统一的业务错误码和错误类型。
//
// 错误码格式采用 HTTP 状态码后加两位细分编号（如 40000），
// 便于前端根据 code 做精确判断而非解析错误文本。
package apperror

import "errors"

// Error 是业务错误类型，包含错误码、用户面消息和可选的内部错误。
type Error struct {
	Code    int
	Message string
	Err     error
}

func (e *Error) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	return e.Message
}

func (e *Error) Unwrap() error {
	return e.Err
}

// New 创建一个新的业务错误。
func New(code int, message string) *Error {
	return &Error{Code: code, Message: message}
}

// Wrap 创建一个携带内部错误的业务错误，用于保留底层错误链。
func Wrap(code int, message string, err error) *Error {
	return &Error{Code: code, Message: message, Err: err}
}

// As 从 error 链中提取 *Error，方便调用方判断业务错误类型。
func As(err error) (*Error, bool) {
	var appErr *Error
	if errors.As(err, &appErr) {
		return appErr, true
	}
	return nil, false
}

const (
	CodeOK              = 0      // 成功
	CodeInvalidArgument = 40000  // 请求参数校验失败
	CodeUnauthorized    = 40100  // 未认证或 Token 无效
	CodeForbidden       = 40300  // 无权限访问
	CodeNotFound        = 40400  // 请求资源不存在
	CodeConflict        = 40900  // 资源冲突（如用户名已存在）
	CodeTooManyRequests = 42900  // 请求频率超限
	CodeInternal        = 50000  // 服务器内部错误
)

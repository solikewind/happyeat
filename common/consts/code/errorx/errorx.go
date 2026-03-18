package errorx

import "strings"

const defaultCode errCode = 1001 // defaultCode 用于动态生成错误，code 1001-1999

type CodeError struct {
	// 重点1: 将int改成errCode类型
	Code errCode `json:"code"`
	Msg  string  `json:"msg"`
}

// 1. 静态错误：走 stringer 自动生成的中文
func NewCodeError(code errCode) error {
	return &CodeError{Code: code, Msg: code.String()}
}

// 2. 动态通用错误：固定 1001 码，自定义文字
func NewDefaultError(msg string) error {
	return &CodeError{Code: defaultCode, Msg: msg}
}

// 3. 动态业务错误（新增）：指定业务码 + 自定义文字
func NewCodeErrorWithMsg(code errCode, msg string) error {
	return &CodeError{Code: code, Msg: msg}
}

func (e *CodeError) Error() string {
	return e.Msg
}

func (e *CodeError) Data() any {
	return e
}

// 重点2: 不根据map来校验是否是自定义类型
// 判断错误码是否是我们自定义的错误类型
// 这里我们根据stringer,会为那些没有定义的code生成默认的msg来判断是否是我们自定义的code
func IsCodeError(code int) bool {
	return !strings.HasPrefix(errCode(code).String(), "errCode(")
}

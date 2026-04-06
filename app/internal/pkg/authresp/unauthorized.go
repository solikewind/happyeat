package authresp

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/golang-jwt/jwt/v4"
)

// UnauthorizedCallback 供 rest.WithUnauthorizedCallback 使用，返回与 Casbin 中间件一致的 JSON 信封。
func UnauthorizedCallback(w http.ResponseWriter, _ *http.Request, err error) {
	msg := "未授权或令牌无效"
	var ve *jwt.ValidationError
	if errors.As(err, &ve) {
		switch {
		case ve.Errors&jwt.ValidationErrorExpired != 0:
			msg = "登录已过期，请重新登录"
		case ve.Errors&jwt.ValidationErrorMalformed != 0,
			ve.Errors&jwt.ValidationErrorSignatureInvalid != 0:
			msg = "令牌无效"
		case ve.Errors&jwt.ValidationErrorNotValidYet != 0:
			msg = "令牌尚未生效"
		}
	} else if err != nil {
		msg = "未授权，请重新登录"
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusUnauthorized)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"code": 401,
		"msg":  msg,
	})
}

// Package routenorm 与 Casbin 中间件使用相同的路径归一化规则，供 routecheck 等工具复用。
package routenorm

import (
	"regexp"
	"strings"
)

var idPathSegment = regexp.MustCompile(`/[0-9a-fA-F-]{6,}`)

// NormalizePath 将请求路径转为 Casbin obj：UUID 段替换为 /:id，并去掉尾部斜杠。
func NormalizePath(path string) string {
	normalized := idPathSegment.ReplaceAllString(path, "/:id")
	return strings.TrimRight(normalized, "/")
}

// CanonicalObj 在 NormalizePath 基础上，将任意命名路径参数（如 :role、:id）统一为 :id，便于路由模板与 permission 规则对比。
func CanonicalObj(path string) string {
	p := NormalizePath(path)
	param := regexp.MustCompile(`:[a-zA-Z_][a-zA-Z0-9_]*`)
	return param.ReplaceAllString(p, ":id")
}

package middleware

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strings"

	"github.com/solikewind/happyeat/app/internal/svc"
)

var idPathSegment = regexp.MustCompile(`/[0-9a-fA-F-]{6,}`)

type CasbinMiddleware struct {
	svcCtx *svc.ServiceContext
}

func NewCasbinMiddleware(svcCtx *svc.ServiceContext) *CasbinMiddleware {
	return &CasbinMiddleware{svcCtx: svcCtx}
}

func (m *CasbinMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 跳过公开接口；其余接口默认拒绝，需显式策略放行。
		if isPublicPath(r.URL.Path) {
			next(w, r)
			return
		}

		sub := extractSubject(r)
		if sub == "" {
			writeForbidden(w, "forbidden: missing subject")
			return
		}

		obj := normalizePath(r.URL.Path)
		act := strings.ToUpper(r.Method)
		ok, err := m.svcCtx.Casbin.Enforcer.Enforce(sub, obj, act)
		if err != nil {
			writeForbidden(w, "forbidden: casbin error")
			return
		}
		if !ok {
			writeForbidden(w, "forbidden: insufficient permissions")
			return
		}

		next(w, r)
	}
}

func isPublicPath(path string) bool {
	switch path {
	case "/health", "/openapi/happyeat.json", "/central/v1/auth/login":
		return true
	default:
		return false
	}
}

func normalizePath(path string) string {
	normalized := idPathSegment.ReplaceAllString(path, "/:id")
	return strings.TrimRight(normalized, "/")
}

func extractSubject(r *http.Request) string {
	if userClaims, ok := r.Context().Value("user").(map[string]any); ok {
		if role, ok := userClaims["role"].(string); ok && strings.TrimSpace(role) != "" {
			return strings.TrimSpace(role)
		}
		if sub, ok := userClaims["sub"].(string); ok && strings.TrimSpace(sub) != "" {
			return strings.TrimSpace(sub)
		}
	}

	if role := strings.TrimSpace(r.Header.Get("X-Role")); role != "" {
		return role
	}
	return ""
}

func writeForbidden(w http.ResponseWriter, msg string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusForbidden)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"code": 403,
		"msg":  msg,
	})
}

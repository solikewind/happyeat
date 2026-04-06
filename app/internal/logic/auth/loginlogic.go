// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package auth

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/solikewind/happyeat/app/internal/svc"
	"github.com/solikewind/happyeat/app/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

// 开发环境默认账号（仅用于 Swagger/联调，生产应改为真实登录或配置）
const devUsername = "admin"
const devPassword = "admin123"

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 登录获取 JWT，用于 Swagger Authorize 或前端
func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginReq) (*types.LoginReply, error) {
	if req.Username != devUsername || req.Password != devPassword {
		return nil, errors.New("用户名或密码错误")
	}
	const subject = "dev-admin"
	// 与 seedDefaultMappings 中 dev-admin 绑定的一致；前端 AuthContext 解析 JWT 的 role 做菜单/路由权限
	const primaryRole = "super_admin"
	if err := l.svcCtx.Rbac.EnsureUser(subject); err != nil {
		return nil, err
	}

	secret := l.svcCtx.Config.Auth.AccessSecret
	expire := l.svcCtx.Config.Auth.AccessExpire
	if expire <= 0 {
		expire = 86400
	}
	iat := time.Now().Unix()
	exp := iat + expire

	// user_code：Casbin 主体；role：非标准 claim，供前端解析（go-zero 会写入 context，与 Casbin 的 g 策略无关）
	claims := jwt.MapClaims{
		"exp":        exp,
		"iat":        iat,
		"sub":        subject,
		"user_code": subject,
		"role":       primaryRole,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(secret))
	if err != nil {
		return nil, err
	}

	return &types.LoginReply{
		AccessToken: tokenStr,
		Expire:      exp,
	}, nil
}

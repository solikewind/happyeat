package errorx

//go:generate stringer -linecomment -type errCode

type errCode int

// 通用错误码
const (
	ERR_OK errCode = 200 // 正确
)

// 用户服务错误码,200000-299999
const (
	ERR_USER_NOT_EXISTS errCode = 200001 // 用户不存在
)

// 菜单服务错误码,300000-399999
const (
	ERR_MENU_NOT_EXISTS errCode = 300001 // 菜单不存在
)

// 其他服务错误码...

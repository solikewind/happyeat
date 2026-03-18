package errorx

import "github.com/solikewind/happyeat/dal/model/ent"

func FromEnt(err error, code errCode) error {
	if err == nil {
		return nil
	}
	if ent.IsNotFound(err) {
		return NewCodeError(code)
	}
	// 其他数据库错误（如约束冲突、连接失败）统一转为默认错误
	return NewDefaultError(err.Error())
}

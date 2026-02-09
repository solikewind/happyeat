package svc

import "github.com/casbin/casbin/v2"

type CasbinEnforcer struct {
	Enforcer *casbin.Enforcer
}

func NewCasbinEnforcer(modelPath, config string) (*CasbinEnforcer, error) {
	// 初始化 Casbin Enforcer
	return nil, nil
}

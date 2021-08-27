package controllers

import (
	"os"

	"github.com/casbin/casbin"
)

type Provider struct {
	rules *casbin.Enforcer
}

func (m *Provider) SetRBAC(model, policy string) error {
	path, err := os.Getwd()

	if err != nil {
		return err
	}

	m.rules, err = casbin.NewEnforcerSafe(path+model, path+policy, false)

	return err
}

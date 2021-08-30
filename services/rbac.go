package services

import (
	"os"

	"github.com/casbin/casbin"
)

var (
	RBACrules *casbin.Enforcer
)

func SetRBAC(model, policy string) error {
	path, err := os.Getwd()

	if err != nil {
		return err
	}

	RBACrules, err = casbin.NewEnforcerSafe(path+model, path+policy, false)

	return err
}

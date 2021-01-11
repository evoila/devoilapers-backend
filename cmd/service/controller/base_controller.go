package controller

import (
	"OperatorAutomation/cmd/service/user"
	"OperatorAutomation/pkg/core"
)

type BaseController struct {
	Core *core.Core
	UserManagement user.UserManagement
}

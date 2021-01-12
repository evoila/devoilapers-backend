package controller

import (
	"OperatorAutomation/cmd/service/management"
	"OperatorAutomation/pkg/core"
)

// Base controller which provides access to the core and to the user-management.
// This is useful for the web-request handler
type BaseController struct {
	Core           *core.Core
	UserManagement user.UserManagement
}

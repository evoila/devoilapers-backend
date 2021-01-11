package core

import "OperatorAutomation/pkg/core/users"

// Use CreateCore to create an instance.
// Holds all references for interacting with the system from a higher level
type Core struct {
	UserContextManagement *users.UserContextManagement
}

// Creates an instance of the core struct holding the references for
// the user-context-management and services
func CreateCore(UserContextManagement *users.UserContextManagement) *Core {
	return &Core{UserContextManagement: UserContextManagement}
}
package user

import (
	"OperatorAutomation/cmd/service/config"
)

type UserManagement struct {
	users map[string]config.User
}

func CreateUserManagement(users []config.User) UserManagement {
	userManagement := UserManagement{
		users: map[string]config.User{},
	}

	// Loop all users
	for _, userInformation := range users {
		// Create user-context objects which combine the given user information with advanced
		// functionally that requires kubernetes access
		userManagement.users[userInformation.Name] = userInformation
	}

	return userManagement
}

// Delivers user information by username and password
func (ctx UserManagement) TryGetUserInformation(username string, password string) (config.User, bool) {
	user, userCouldBeFound := ctx.users[username]
	if !userCouldBeFound || user.Password != password {
		return user, false
	}

	return user, true
}


// Delivers user information by username and password
func (ctx UserManagement) GetUserInformation(username string, password string) config.User {
	user, userCouldBeFound := ctx.TryGetUserInformation(username, password)
	if userCouldBeFound {
		return user
	}

	panic("User not found")
}

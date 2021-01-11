package users

import "OperatorAutomation/pkg/kubernetes"

// Holds user contexts
type UserContextManagement struct {
	Users map[string]UserContext
}

// Creates an instance of the UserContextManagement struct
func CreateUserContextManagement(userInformations []IUserInformation) UserContextManagement {
	userManagement := UserContextManagement{
		Users: map[string]UserContext{},
	}

	// Loop all users
	for _, userInformation := range userInformations{
		// Create user-context objects which combine the given user information with advanced
		// functionally that requires kubernetes access
		userManagement.Users[userInformation.GetName()] =  UserContext{
			userInformation,
			kubernetes.CreateKubernetesWrapper(userInformation),
		}
 	}

 	return userManagement
}

// Delivers user information by username and password
func (ctx UserContextManagement) GetUserInformation(username string, password string) (*IUserInformation, bool) {
	user, userCouldBeFound := ctx.Users[username]
	if !userCouldBeFound {
		return &user.IUserInformation, false
	}

	return &user.IUserInformation, user.GetPassword() == password
}
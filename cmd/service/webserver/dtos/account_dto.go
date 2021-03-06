package dtos

type AccountCredentialsDto struct {
	Username string `json:"username" example:"admin"`
	Password string `json:"password" example:"masterkey"`
}

type AuthenticationResponseDataDto struct {
	IsValid bool   `json:"isValid" example:"true"`
	Role    string `json:"role" example:"admin"`
}

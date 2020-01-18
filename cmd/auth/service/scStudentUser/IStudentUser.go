package scStudentUser

type ScStudentUser struct {
	IScStudentUser
}

func New() ScStudentUser {
	scStudentUser := &StudentuserV1{}
	scStudentUser.Init()
	return ScStudentUser{scStudentUser}
}

type IScStudentUser interface {
	Login(*LoginRequest) (*LoginResult, error)
}
type LoginRequest struct {
	Mobile  string
	SmsCode string
}
type TokenResult struct {
	Token     string `json:"token"`
	TokenType string `json:"token_type"`
	ExpiresIn int32  `json:"expires_in"`
}
type LoginResult struct {
	Token *TokenResult `json:"token"`
}

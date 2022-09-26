package internal

var (
	RequestStorage  = make([]AuthenticateRequest, 0)
	ClientStorage   = make([]Client, 0)
	UserStorage     = make([]User, 0)
	AuthCodeStorage = make([]AuthCode, 0)
	TokenStorage    = make([]Token, 0)
)

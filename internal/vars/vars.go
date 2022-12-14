package vars

var (
	// openid provider
	// resource server
	RS_HOST = "localhost:7000"
	RS_URL  = "http://localhost:7000"
	OP_HOST = "localhost:7001"
	OP_URL  = "http://localhost:7001"
	RP_HOST = "localhost:7002"
	RP_URL  = "http://localhost:7002"

	REDIRECT_URI = "/callback"

	CHECK_TOKEN_ENDPOINT  = "/check-token"
	RS_PROTECTED_ENDPOINT = "/protected"
)

var (
	AUTHN_ENDPOINT = "/authenticate"
	LOGIN_ENDPOINT = "/login"
	TOKEN_ENDPOINT = "/token"
)

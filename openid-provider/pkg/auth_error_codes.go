package pkg

// this are error codes in callback response
const (
	InvalidRequest          string = "invalid_request"
	UnauthorizedClient      string = "unauthorized_client"
	AccessDenied            string = "access_denied"
	UnsupportedResponseType string = "unsupported_response_type"
	InvalidScope            string = "invalid_scope"
	ServerError             string = "server_error"
	TemporaryUnavailable    string = "temporary_unavailable"
)

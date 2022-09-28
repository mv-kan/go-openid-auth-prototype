package internal

// type AuthCode string
type AuthCode struct {
	ID            string
	AuthRequestID string
	Scope         []string
	ClientID      string
	RedirectURI   string
	State         string
}

func (c AuthCode) GetID() string {
	return c.ID
}

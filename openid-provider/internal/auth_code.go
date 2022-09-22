package internal

// type AuthCode string
type AuthCode struct {
	ID            string
	AuthRequestID string
}

func (c AuthCode) GetID() string {
	return c.ID
}

package internal

type Token struct {
	Scopes      []string
	ID          string
	ClientID    string
	AccessToken string
}

func (t Token) GetID() string {
	return t.ID
}

package internal

type Client struct {
	ID           string
	Secret       string
	RedirectURIs []string
}

func (c Client) GetID() string {
	return c.ID
}

package fb

import "context"

// Client is a new firebase client.
type Client struct {
	token string
}

func NewClient(token string) *Client {
	return &Client{
		token: token,
	}
}

// Users grabs all the associated users from the database.
func (c *Client) Users(ctx context.Context) ([]string, error) {
	// TODO: Implement this.
	return nil, nil
}

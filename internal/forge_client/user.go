package forge_client

import "context"

type User struct {
	ID        int64   `json:"id"`
	Name      string  `json:"name"`
	Email     string  `json:"email"`
	CreatedAt *string `json:"created_at"`
	UpdatedAt *string `json:"updated_at"`
}

func (c *Client) GetUser(ctx context.Context) (*User, error) {
	var user User
	id, err := c.GetJsonApi(ctx, "/user", &user)
	if err != nil {
		return nil, err
	}
	user.ID = int64(id)
	return &user, nil
}

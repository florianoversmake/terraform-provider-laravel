package forge_client

import "context"

type User struct {
	ID                      int64  `json:"id"`
	Name                    string `json:"name"`
	Email                   string `json:"email"`
	CardLastFour            string `json:"card_last_four"`
	ConnectedToGithub       bool   `json:"connected_to_github"`
	ConnectedToGitlab       bool   `json:"connected_to_gitlab"`
	ConnectedToBitbucketTwo bool   `json:"connected_to_bitbucket_two"`
	ConnectedToDigitalocean bool   `json:"connected_to_digitalocean"`
	ConnectedToLinode       bool   `json:"connected_to_linode"`
	ConnectedToVultr        bool   `json:"connected_to_vultr"`
	ConnectedToAWS          bool   `json:"connected_to_aws"`
	ConnectedToHetzner      bool   `json:"connected_to_hetzner"`
	ReadyForBilling         bool   `json:"ready_for_billing"`
	StripeIsActive          int    `json:"stripe_is_active"`
	StripePrice             string `json:"stripe_price"`
	Subscribed              int    `json:"subscribed"`
	CanCreateServers        bool   `json:"can_create_servers"`
	TwoFAEnabled            bool   `json:"2fa_enabled"`
}

func (c *Client) GetUser(ctx context.Context) (*User, error) {
	var res struct {
		User User `json:"user"`
	}
	if err := c.doRequest(ctx, "GET", "/user", nil, &res); err != nil {
		return nil, err
	}
	return &res.User, nil
}

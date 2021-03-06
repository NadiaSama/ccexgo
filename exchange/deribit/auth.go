package deribit

import (
	"context"
	"time"
)

type (
	AuthParam struct {
		GrantType    string `json:"grant_type"`
		ClientID     string `json:"client_id"`
		ClientSecret string `json:"client_secret"`
	}

	AuthResult struct {
		AccessToken  string `json:"access_token"`
		ExpiresIn    int    `json:"expires_in"`
		RefreshToken string `json:"refresh_token"`
		Scope        string `json:"scope"`
		TokenType    string `json:"token_type"`
	}

	Token interface {
		SetToken(string)
	}

	AuthToken struct {
		AccessToken string `json:"access_token"`
	}
)

func (ap *AuthToken) SetToken(token string) {
	ap.AccessToken = token
}

func (c *Client) getToken(ctx context.Context) (string, error) {
	now := time.Now()
	if c.expire.After(now) {
		return c.accessToken, nil
	}

	var r AuthResult
	param := &AuthParam{
		ClientID:     c.key,
		ClientSecret: c.secret,
		GrantType:    "client_credentials",
	}
	if err := c.call(ctx, "public/auth", param, &r, false); err != nil {
		return "", err
	}
	c.tokenMu.Lock()
	defer c.tokenMu.Unlock()
	c.accessToken = r.AccessToken
	c.expire = now.Add(time.Duration(r.ExpiresIn-1) * time.Second)
	return r.AccessToken, nil
}

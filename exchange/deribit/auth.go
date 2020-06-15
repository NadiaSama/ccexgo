package deribit

import (
	"time"
)

type (
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

	AuthParam struct {
		AccessToken string `json:"access_token"`
	}
)

func (ap *AuthParam) SetToken(token string) {
	ap.AccessToken = token
}

func (c *Client) getToken() (string, error) {
	now := time.Now()
	if c.expire.After(now) {
		return c.accessToken, nil
	}

	var r AuthResult
	if err := c.Client.Conn.Call(c.Ctx, "public/auth", nil, &r); err != nil {
		return "", err
	}
	c.tokenMu.Lock()
	defer c.tokenMu.Unlock()
	c.accessToken = r.AccessToken
	c.expire = now.Add(time.Duration(r.ExpiresIn) * time.Second)
	return r.AccessToken, nil
}

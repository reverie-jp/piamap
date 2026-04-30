package google

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const googleUserInfoURL = "https://www.googleapis.com/oauth2/v2/userinfo"

type UserInfo struct {
	Sub     string `json:"id"`
	Email   string `json:"email"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
}

type AuthClient struct {
	config *oauth2.Config
}

func NewAuthClient(clientID, clientSecret, redirectURL string) *AuthClient {
	return &AuthClient{
		config: &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  redirectURL,
			Scopes:       []string{"openid", "profile", "email"},
			Endpoint:     google.Endpoint,
		},
	}
}

// Exchange は OAuth code を Google で検証し、profile + email を取得する。
func (c *AuthClient) Exchange(ctx context.Context, code string) (*UserInfo, error) {
	token, err := c.config.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code: %w", err)
	}

	client := c.config.Client(ctx, token)
	resp, err := client.Get(googleUserInfoURL)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("google userinfo returned status %d", resp.StatusCode)
	}

	var userInfo UserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, fmt.Errorf("failed to decode user info: %w", err)
	}
	return &userInfo, nil
}

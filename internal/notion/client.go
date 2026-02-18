package notion

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	baseURL        = "https://api.notion.com/v1"
	notionVersion  = "2022-06-28"
)

// Client calls the Notion REST API directly.
type Client struct {
	token      string
	httpClient *http.Client
}

// NewClient creates a Client authenticated with the given integration token.
func NewClient(token string) *Client {
	return &Client{
		token:      token,
		httpClient: &http.Client{},
	}
}

// do executes an authenticated request and decodes the JSON response into dst.
func (c *Client) do(ctx context.Context, method, path string, dst any) error {
	req, err := http.NewRequestWithContext(ctx, method, baseURL+path, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Notion-Version", notionVersion)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var apiErr struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&apiErr); err != nil {
			return fmt.Errorf("HTTP %d", resp.StatusCode)
		}
		return fmt.Errorf("%s: %s", apiErr.Code, apiErr.Message)
	}

	return json.NewDecoder(resp.Body).Decode(dst)
}

// Ping verifies the token is valid by fetching the bot user via GET /v1/users/me.
// Returns the bot's display name on success.
func (c *Client) Ping(ctx context.Context) (string, error) {
	var user struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}
	if err := c.do(ctx, http.MethodGet, "/users/me", &user); err != nil {
		return "", err
	}

	name := user.Name
	if name == "" {
		name = user.ID
	}
	return name, nil
}

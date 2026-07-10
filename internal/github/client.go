package github

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/myerscode/aws-meta/internal/util"
)

type Client struct {
	client *http.Client
	token  string
}

func NewGitHubClient(token string) Client {
	return Client{
		client: &http.Client{Timeout: 10 * time.Second},
		token:  token,
	}
}

func (c Client) Fetch(url string) ([]byte, error) {

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/vnd.github.v3+json")
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	} else if token := os.Getenv("AWSMETA_GITHUB_TOKEN"); token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	} else {
		// No token provided, proceed without authentication
		// This may result in rate limiting for unauthenticated requests
		util.LogWarning("No GitHub token provided. API requests may be rate-limited.")
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	err = resp.Body.Close()

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API error (%d): %s", resp.StatusCode, string(body))
	}

	return body, nil
}

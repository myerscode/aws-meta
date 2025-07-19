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
	client     *http.Client
	token      string
	rateLimit  int
	rateRemain int
	rateReset  time.Time
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
	} else if token := os.Getenv("GITHUB_TOKEN"); token != "" {
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

	c.updateRateLimit(resp)

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

// updateRateLimit updates the rate limit information based on the response headers
func (c Client) updateRateLimit(resp *http.Response) {
	if limit := resp.Header.Get("X-RateLimit-Limit"); limit != "" {
		_, err := fmt.Sscanf(limit, "%d", &c.rateLimit)
		if err != nil {
			util.LogError(fmt.Sprintf("Error parsing rate limit: %v", err))
		}
	}

	if remain := resp.Header.Get("X-RateLimit-Remaining"); remain != "" {
		_, err := fmt.Sscanf(remain, "%d", &c.rateRemain)
		if err != nil {
			util.LogError(fmt.Sprintf("Error parsing rate remaining: %v", err))
		}
	}

	if reset := resp.Header.Get("X-RateLimit-Reset"); reset != "" {
		var resetUnix int64
		_, err := fmt.Sscanf(reset, "%d", &resetUnix)
		if err != nil {
			util.LogError(fmt.Sprintf("Error parsing rate reset: %v", err))
		} else {
			c.rateReset = time.Unix(resetUnix, 0)
		}
	}
}

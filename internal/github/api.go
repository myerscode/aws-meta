package github

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/myerscode/aws-meta/internal/util"
	"github.com/pterm/pterm"

	"github.com/pkg/errors"
)

type Config struct {
	Owner    string
	RepoName string
	Branch   string
}

type Repo struct {
	Config
	Client
}

type RepoTag struct {
	Name       string `json:"name"`
	Commit     Commit `json:"commit"`
	ZipBallURL string `json:"zip_ball_url"`
	TarBallURL string `json:"tarball_url"`
}
type Commit struct {
	SHA string `json:"sha"`
	URL string `json:"url"`
}

type RepoTreeItem struct {
	Path string `json:"path"`
	Url  string `json:"url"`
	Sha  string `json:"sha"`
}

func (r Repo) LatestTag() (RepoTag, error) {
	path := fmt.Sprintf("repos/%s/%s/tags", r.Owner, r.RepoName)

	data, err := githubAPIWithGetMethod(path)

	if err != nil {
		return RepoTag{}, err
	}

	var tags []RepoTag

	if err := json.Unmarshal(data, &tags); err != nil {
		return RepoTag{}, err
	}

	return tags[0], nil
}

func (r Repo) GetGithubRepoTrees(commitSha string, directory string) ([]RepoTreeItem, error) {

	path := fmt.Sprintf("repos/%s/%s/git/trees/%s:%s?recursive=1", r.Owner, r.RepoName, commitSha, url.QueryEscape(directory))

	githubPath := fmt.Sprintf("https://api.github.com/%s", path)

	data, err := fetchContent(githubPath)

	if err != nil {
		return nil, errors.Wrapf(err, "fail to get trees (%s)", path)
	}

	tree := struct {
		Tree []RepoTreeItem `json:"tree"`
	}{}

	if err := json.Unmarshal(data, &tree); err != nil {
		return nil, err
	}

	return tree.Tree, nil
}

func (r Repo) GetBlobFromTag(tagName string, filename string) ([]byte, error) {
	apiUrl := fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/refs/tags/%s/%s/%s", r.Owner, r.RepoName, tagName, r.RepoName, filename)

	return fetchContent(apiUrl)
}

func (r Repo) FetchTags(perPage int) ([]RepoTag, error) {
	apiUrl := fmt.Sprintf("https://api.github.com/repos/%s/%s/tags?per_page=%d", r.Owner, r.RepoName, perPage)

	blob, err := fetchContent(apiUrl)

	if err != nil {
		return nil, err
	}

	var tags []RepoTag

	if err = json.Unmarshal(blob, &tags); err != nil {
		return []RepoTag{}, err
	}

	return tags, nil
}

func fetchContent(apiUrl string) ([]byte, error) {

	util.LogTrace(fmt.Sprintf("Getting: %s", apiUrl))

	resp, err := http.Get(apiUrl)

	if err != nil {
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		closeErr := Body.Close()
		if closeErr != nil {
			util.LogError("Failed to close response body", []pterm.LoggerArgument{{Key: "error", Value: closeErr}})
		}
	}(resp.Body)

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("GitHub API error (%d / %s): %s", resp.StatusCode, resp.Status, string(body))
	}

	return data, nil
}

func GetGithubRepoBlobs(owner string, repo string, branchOrSha string, filename string) ([]byte, error) {

	apiUrl := fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/%s/%s", owner, repo, branchOrSha, filename)

	util.LogTrace(fmt.Sprintf("Getting: %s", apiUrl))

	resp, err := http.Get(apiUrl)

	if err != nil {
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		closeErr := Body.Close()
		if closeErr != nil {
			util.LogError("Failed to close response body", []pterm.LoggerArgument{{Key: "error", Value: closeErr}})
		}
	}(resp.Body)

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("[%d %s] %s", resp.StatusCode, resp.Status, string(data[:]))
	}

	return data, nil
}

func githubAPIWithGetMethod(path string) ([]byte, error) {

	httpClient := http.Client{
		Timeout: time.Second * 10,
	}

	apiUrl := fmt.Sprintf("https://api.github.com/%s", path)

	req, err := http.NewRequest("GET", apiUrl, nil)
	if err != nil {
		return nil, err
	}

	// Add common headers
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	// Check for GITHUB_TOKEN and add authorization if present
	if token := os.Getenv("GITHUB_TOKEN"); token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	} else {
		// No token provided, proceed without authentication
		// This may result in rate limiting for unauthenticated requests
		util.LogError("No GitHub token provided. API requests may be rate-limited.", nil)
	}

	fmt.Printf("GET %s\n", apiUrl)

	resp, err := httpClient.Do(req)

	if err != nil {
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		closeErr := Body.Close()
		if closeErr != nil {
			util.LogError("Failed to close response body", []pterm.LoggerArgument{{Key: "error", Value: closeErr}})
		}
	}(resp.Body)

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("%s: %s", resp.Status, string(data[:]))
	}

	return data, nil
}

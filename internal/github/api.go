package github

import (
	"encoding/json"
	"fmt"
)

type Config struct {
	Owner    string
	RepoName string
}

type Repo struct {
	Config
	Client
}

type RepoTag struct {
	Name string `json:"name"`
}

func (r Repo) FetchTags(perPage int) ([]RepoTag, error) {
	apiUrl := fmt.Sprintf("https://api.github.com/repos/%s/%s/tags?per_page=%d", r.Owner, r.RepoName, perPage)

	blob, err := r.Fetch(apiUrl)

	if err != nil {
		return nil, err
	}

	var tags []RepoTag

	if err = json.Unmarshal(blob, &tags); err != nil {
		return []RepoTag{}, err
	}

	return tags, nil
}

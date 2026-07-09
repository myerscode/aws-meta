package github

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/myerscode/aws-meta/internal/util"
)

// TarballFiles holds extracted file contents keyed by their relative path
// within the repository (e.g. "botocore/data/partitions.json").
type TarballFiles map[string][]byte

// DownloadAndExtract fetches the tarball for a given tag and extracts only
// files whose paths (relative to the repo root) match any of the provided
// prefix filters. This replaces hundreds of individual API calls with a
// single HTTP request for the compressed archive.
func (r Repo) DownloadAndExtract(tag RepoTag, pathPrefixes []string) (TarballFiles, error) {
	tarballURL := fmt.Sprintf("https://github.com/%s/%s/archive/refs/tags/%s.tar.gz", r.Owner, r.RepoName, tag.Name)

	util.LogInfo(fmt.Sprintf("Downloading tarball for tag %s", tag.Name))

	req, err := http.NewRequest("GET", tarballURL, nil)
	if err != nil {
		return nil, fmt.Errorf("creating tarball request: %w", err)
	}

	// Use auth token if available to avoid rate limiting
	if r.Client.token != "" {
		req.Header.Set("Authorization", "Bearer "+r.Client.token)
	} else if token := os.Getenv("AWSMETA_GITHUB_TOKEN"); token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	} else {
		util.LogWarning("No AWSMETA_GITHUB_TOKEN set. Tarball download may be rate-limited.")
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("downloading tarball: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("tarball download failed with status %d", resp.StatusCode)
	}

	return extractFromTarGz(resp.Body, pathPrefixes)
}

// extractFromTarGz streams a gzipped tar archive and extracts files matching
// any of the given path prefixes. Paths are matched after stripping the
// top-level directory that GitHub adds (e.g. "botocore-1.35.100/").
func extractFromTarGz(reader io.Reader, pathPrefixes []string) (TarballFiles, error) {
	gz, err := gzip.NewReader(reader)
	if err != nil {
		return nil, fmt.Errorf("creating gzip reader: %w", err)
	}
	defer gz.Close()

	tr := tar.NewReader(gz)
	files := make(TarballFiles)

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("reading tar entry: %w", err)
		}

		// Skip directories
		if header.Typeflag != tar.TypeReg {
			continue
		}

		// Strip the top-level directory (e.g. "botocore-1.35.100/botocore/data/..." -> "botocore/data/...")
		relativePath := stripTopDir(header.Name)

		if !matchesAnyPrefix(relativePath, pathPrefixes) {
			continue
		}

		data, err := io.ReadAll(tr)
		if err != nil {
			return nil, fmt.Errorf("reading file %s from tar: %w", relativePath, err)
		}

		files[relativePath] = data
	}

	util.LogInfo(fmt.Sprintf("Extracted %d files from tarball", len(files)))

	return files, nil
}

// stripTopDir removes the first path component from a tar entry path.
// GitHub tarballs always have a top-level directory like "reponame-tagname/".
func stripTopDir(path string) string {
	idx := strings.IndexByte(path, '/')
	if idx < 0 {
		return path
	}
	return path[idx+1:]
}

// matchesAnyPrefix returns true if the path starts with any of the given prefixes.
func matchesAnyPrefix(path string, prefixes []string) bool {
	for _, prefix := range prefixes {
		if strings.HasPrefix(path, prefix) {
			return true
		}
	}
	return false
}

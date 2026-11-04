package documents

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type githubRepository struct {
	DefaultBranch string `json:"default_branch"`
	Private       bool   `json:"private"`
}

type githubCommit struct {
	SHA string `json:"sha"`
}

func fetchGitHubArchive(ctx context.Context, repositoryURL, sourceDir, token string) (string, string, string, error) {
	parsed, err := url.Parse(repositoryURL)
	if err != nil || parsed.Scheme != "https" || parsed.Host != "github.com" {
		return "", "", "", invalidSource("invalid_github_url", "The GitHub URL is invalid.")
	}
	parts := strings.Split(strings.Trim(parsed.Path, "/"), "/")
	if len(parts) != 2 {
		return "", "", "", invalidSource("invalid_github_url", "The GitHub URL is invalid.")
	}

	client := &http.Client{Timeout: 45 * time.Second, CheckRedirect: githubRedirect}
	base := "https://api.github.com/repos/" + url.PathEscape(parts[0]) + "/" + url.PathEscape(parts[1])

	var repository githubRepository
	if err := githubJSON(ctx, client, base, token, &repository); err != nil {
		return "", "", "", err
	}
	if repository.Private || repository.DefaultBranch == "" {
		return "", "", "", invalidSource("github_unavailable", "Use a public repository with a default branch.")
	}

	var commit githubCommit
	commitURL := base + "/commits/" + url.PathEscape(repository.DefaultBranch)
	if err := githubJSON(ctx, client, commitURL, token, &commit); err != nil {
		return "", "", "", err
	}
	if len(commit.SHA) != 40 || !isHex(commit.SHA) {
		return "", "", "", invalidSource("github_response", "GitHub returned an invalid commit ID.")
	}

	archiveURL := base + "/zipball/" + commit.SHA
	archivePath, digest, err := githubArchive(ctx, client, archiveURL, sourceDir, token)
	if err != nil {
		return "", "", "", err
	}

	return archivePath, commit.SHA, digest, nil
}

func githubRedirect(req *http.Request, via []*http.Request) error {
	if len(via) >= 4 || req.URL.Scheme != "https" || req.URL.User != nil || req.URL.Port() != "" {
		return fmt.Errorf("unsafe GitHub redirect")
	}
	if req.URL.Hostname() != "api.github.com" && req.URL.Hostname() != "codeload.github.com" {
		return fmt.Errorf("unsafe GitHub redirect")
	}
	if req.URL.Hostname() == "codeload.github.com" {
		req.Header.Del("Authorization")
	}

	return nil
}

func githubJSON(ctx context.Context, client *http.Client, endpoint, token string, result any) error {
	response, err := githubGet(ctx, client, endpoint, token)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if err := githubStatus(response); err != nil {
		return err
	}
	data, err := io.ReadAll(io.LimitReader(response.Body, 1<<20+1))
	if err != nil {
		return fmt.Errorf("read GitHub response: %w", err)
	}
	if len(data) > 1<<20 || json.Unmarshal(data, result) != nil {
		return invalidSource("github_response", "GitHub returned an invalid response.")
	}

	return nil
}

func githubArchive(ctx context.Context, client *http.Client, endpoint, sourceDir, token string) (string, string, error) {
	response, err := githubGet(ctx, client, endpoint, token)
	if err != nil {
		return "", "", err
	}
	defer response.Body.Close()

	if err := githubStatus(response); err != nil {
		return "", "", err
	}
	if response.ContentLength > MaxArchiveBytes {
		return "", "", invalidSource("archive_too_large", "The repository archive is too large.")
	}

	file, err := os.CreateTemp(sourceDir, "github-*.zip")
	if err != nil {
		return "", "", fmt.Errorf("create GitHub archive: %w", err)
	}
	keep := false
	defer func() {
		_ = file.Close()
		if !keep {
			_ = os.Remove(file.Name())
		}
	}()

	hash := sha256.New()
	size, err := io.Copy(io.MultiWriter(file, hash), io.LimitReader(&contextReader{ctx: ctx, reader: response.Body}, MaxArchiveBytes+1))
	if err != nil {
		return "", "", fmt.Errorf("download GitHub archive: %w", err)
	}
	if size > MaxArchiveBytes {
		return "", "", invalidSource("archive_too_large", "The repository archive is too large.")
	}
	if err := file.Close(); err != nil {
		return "", "", fmt.Errorf("close GitHub archive: %w", err)
	}

	keep = true

	return file.Name(), hex.EncodeToString(hash.Sum(nil)), nil
}

func githubGet(ctx context.Context, client *http.Client, endpoint, token string) (*http.Response, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("create GitHub request: %w", err)
	}
	request.Header.Set("Accept", "application/vnd.github+json")
	request.Header.Set("User-Agent", "ai-doc-gen")
	request.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	if token != "" {
		request.Header.Set("Authorization", "Bearer "+token)
	}

	response, err := client.Do(request)
	if err != nil {
		return nil, &IntakeError{Code: "github_network", Message: "GitHub is temporarily unavailable.", Temporary: true}
	}

	return response, nil
}

func githubStatus(response *http.Response) error {
	switch response.StatusCode {
	case http.StatusOK:
		return nil
	case http.StatusNotFound:
		return invalidSource("github_not_found", "The repository was not found or is empty.")
	case http.StatusUnauthorized:
		return invalidSource("github_auth", "GitHub access failed. Check the server credentials.")
	case http.StatusForbidden, http.StatusTooManyRequests:
		body, _ := io.ReadAll(io.LimitReader(response.Body, 4096))
		limited := strings.Contains(strings.ToLower(string(body)), "rate limit")
		if limited || response.Header.Get("X-RateLimit-Remaining") == "0" ||
			response.Header.Get("Retry-After") != "" || response.StatusCode == http.StatusTooManyRequests {

			return invalidSource("github_rate_limit", "GitHub rate limit reached. Try again later.")
		}

		return invalidSource("github_forbidden", "GitHub denied access to this repository.")
	default:
		if response.StatusCode >= 500 {
			return &IntakeError{Code: "github_unavailable", Message: "GitHub is temporarily unavailable.", Temporary: true}
		}

		return invalidSource("github_response", "GitHub could not provide this repository.")
	}
}

func isHex(value string) bool {
	for _, char := range value {
		if !strings.ContainsRune("0123456789abcdefABCDEF", char) {
			return false
		}
	}

	return true
}

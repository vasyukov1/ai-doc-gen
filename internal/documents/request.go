package documents

import (
	"errors"
	"net/url"
	"path"
	"strings"
)

type Request struct {
	GitHubURL   string
	HasArchive  bool
	ArchiveName string
	ProfileID   string
	Language    string
	Description string
}

func ValidateRequest(request Request) (Request, error) {
	request.GitHubURL = strings.TrimSpace(request.GitHubURL)
	request.Description = strings.TrimSpace(request.Description)

	if request.GitHubURL == "" && !request.HasArchive {
		return Request{}, errors.New("provide one GitHub URL or ZIP archive")
	}
	if request.GitHubURL != "" && request.HasArchive {
		return Request{}, errors.New("provide one GitHub URL or ZIP archive")
	}
	if request.ProfileID != "academic" && request.ProfileID != "technical-overview" {
		return Request{}, errors.New("choose a document profile")
	}
	if request.Language != "en" && request.Language != "ru" {
		return Request{}, errors.New("choose English or Russian")
	}
	if len(request.Description) > 2000 {
		return Request{}, errors.New("description is too long")
	}

	if request.HasArchive {
		if !strings.EqualFold(path.Ext(request.ArchiveName), ".zip") {
			return Request{}, errors.New("upload a ZIP archive")
		}
		return request, nil
	}

	parsed, err := url.Parse(request.GitHubURL)
	if err != nil || parsed.Scheme != "https" || !strings.EqualFold(parsed.Host, "github.com") {
		return Request{}, errors.New("use a public github.com URL")
	}
	if parsed.User != nil || parsed.RawQuery != "" || parsed.Fragment != "" {
		return Request{}, errors.New("use a repository URL without extra parameters")
	}

	parts := strings.Split(strings.Trim(parsed.Path, "/"), "/")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return Request{}, errors.New("use a repository URL")
	}

	request.GitHubURL = "https://github.com/" + parts[0] + "/" + strings.TrimSuffix(parts[1], ".git")

	return request, nil
}

package documents

import (
	"errors"
	"github.com/vasyukov1/ai-doc-gen/internal/templates"
	"net/url"
	"path"
	"regexp"
	"strings"
)

var githubOwner = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9-]{0,38}$`)
var githubRepo = regexp.MustCompile(`^[A-Za-z0-9._-]{1,100}$`)

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
	profile, err := templates.Load(request.ProfileID)
	if err != nil {
		return Request{}, errors.New("choose a document profile")
	}
	if !profile.SupportsLanguage(request.Language) {
		return Request{}, errors.New("choose English or Russian; Catalan is available for neutral technical profiles only")
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
	if err != nil || parsed.Scheme != "https" || !strings.EqualFold(parsed.Hostname(), "github.com") || parsed.Port() != "" {
		return Request{}, errors.New("use a public github.com URL")
	}
	if parsed.User != nil || parsed.RawQuery != "" || parsed.Fragment != "" || parsed.EscapedPath() != parsed.Path {
		return Request{}, errors.New("use a repository URL without extra parameters")
	}

	parts := strings.Split(strings.Trim(parsed.Path, "/"), "/")
	repo := ""
	if len(parts) == 2 {
		repo = strings.TrimSuffix(parts[1], ".git")
	}
	if len(parts) != 2 || !githubOwner.MatchString(parts[0]) || !githubRepo.MatchString(repo) || repo == "." || repo == ".." {
		return Request{}, errors.New("use a repository URL")
	}

	request.GitHubURL = "https://github.com/" + parts[0] + "/" + repo

	return request, nil
}

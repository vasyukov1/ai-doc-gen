package documents

import (
	"path"
	"strings"
)

var excludedDirectories = map[string]bool{
	".git": true, ".idea": true, ".next": true, ".venv": true,
	"__pycache__": true, "build": true, "coverage": true, "dist": true,
	"node_modules": true, "target": true, "vendor": true,
}

var binaryExtensions = map[string]bool{
	".bmp": true, ".class": true, ".dll": true, ".eot": true,
	".exe": true, ".gif": true, ".ico": true, ".jpeg": true,
	".jpg": true, ".mp3": true, ".mp4": true, ".otf": true,
	".pdf": true, ".png": true, ".so": true, ".ttf": true,
	".webp": true, ".woff": true, ".woff2": true,
}

var archiveExtensions = map[string]bool{
	".7z": true, ".gz": true, ".jar": true, ".rar": true,
	".tar": true, ".tgz": true, ".zip": true,
}

// Lock files and repository metadata describe tooling, not the product.
var noiseNames = map[string]bool{
	".dockerignore": true, ".editorconfig": true, ".gitattributes": true,
	".gitignore": true, ".gitmodules": true, "cargo.lock": true,
	"composer.lock": true, "gemfile.lock": true, "go.sum": true,
	"gradle.lockfile": true, "package-lock.json": true, "pipfile.lock": true,
	"pnpm-lock.yaml": true, "poetry.lock": true, "yarn.lock": true,
}

var sensitiveNames = map[string]bool{
	".git-credentials": true, ".netrc": true, ".npmrc": true,
	".pypirc": true, "auth.json": true, "credentials": true,
	"credentials.json": true, "id_ed25519": true, "id_rsa": true,
	"kubeconfig": true, "service-account.json": true,
	"service_account.json": true,
}

func classifyFile(name, profile string, size uint64) (string, bool, bool) {
	parts := strings.Split(strings.ToLower(name), "/")
	base := parts[len(parts)-1]
	ext := path.Ext(base)

	for _, part := range parts[:len(parts)-1] {
		if part == ".ssh" || part == ".aws" || part == ".kube" ||
			part == "secrets" || part == "credentials" {
			return "likely_secret", false, false
		}
		if excludedDirectories[part] {
			return "dependency_or_output", false, false
		}
		if part == "generated" || part == "gen" {
			return "generated_output", false, false
		}
	}

	if sensitiveNames[base] || base == ".env" || strings.HasPrefix(base, ".env.") ||
		strings.Contains(base, "secret") || ext == ".pem" || ext == ".key" ||
		ext == ".p12" || ext == ".pfx" || ext == ".jks" ||
		ext == ".keystore" || ext == ".tfstate" || ext == ".tfvars" {
		return "likely_secret", false, false
	}
	if noiseNames[base] {
		return "tooling_metadata", false, false
	}
	if archiveExtensions[ext] {
		return "nested_archive", false, false
	}
	if binaryExtensions[ext] {
		return "binary_file", false, false
	}
	if size > maxTextFileBytes {
		return "file_too_large", false, false
	}
	if strings.HasSuffix(base, ".pb.go") || strings.HasSuffix(base, ".min.js") ||
		strings.HasSuffix(base, ".map") || strings.Contains(base, "generated") {
		return "generated_output", false, false
	}

	// Build and deployment files stay available on request but are not part of
	// the default evidence selection.
	if strings.HasPrefix(base, "dockerfile") || strings.HasPrefix(base, "docker-compose.") ||
		base == "compose.yml" || base == "compose.yaml" || base == "makefile" {
		return "build_config", true, profile == "arc42" || profile == "srs"
	}

	if strings.HasSuffix(base, "_test.go") || strings.HasPrefix(base, "test_") ||
		strings.HasSuffix(base, "_test.py") || strings.HasSuffix(base, ".test.ts") ||
		strings.HasSuffix(base, ".spec.ts") || strings.HasSuffix(base, ".test.js") ||
		strings.HasSuffix(base, ".spec.js") ||
		strings.Contains("/"+strings.ToLower(name)+"/", "/tests/") ||
		strings.Contains("/"+strings.ToLower(name)+"/", "/test/") ||
		strings.Contains("/"+strings.ToLower(name)+"/", "/__tests__/") {
		return "test_file", true, profile == "technical-overview" || profile == "srs" || profile == "arc42"
	}

	return "source_file", true, true
}

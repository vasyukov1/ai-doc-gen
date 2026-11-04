package documents

import (
	"archive/zip"
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/vasyukov1/ai-doc-gen/internal/domain"
	"github.com/vasyukov1/ai-doc-gen/internal/templates"
)

var MetadataFields = []string{"project", "organization", "faculty", "program", "author", "group", "supervisor", "supervisor_position", "approver", "approver_position", "year", "code"}

func ValidateDocumentInput(input domain.DocumentInput, profile, language string) (domain.DocumentInput, error) {
	p, err := templates.Load(profile)
	if err != nil {
		return input, err
	}
	if input.Metadata == nil {
		input.Metadata = map[string]string{}
	}
	for _, key := range MetadataFields {
		value := strings.TrimSpace(input.Metadata[key])
		if len(value) > 600 || !safeText(value) {
			return input, fmt.Errorf("title field %s must be valid text below 600 bytes", key)
		}
		input.Metadata[key] = value
	}
	if input.Metadata["project"] == "" {
		return input, errors.New("enter a project title")
	}
	if len(input.Specification) > 64000 || !safeText(input.Specification) {
		return input, errors.New("specification must be UTF-8 text below 64000 bytes")
	}
	if p.Academic && strings.TrimSpace(input.Specification) == "" {
		return input, errors.New("provide the technical specification as text or complete the structured fields")
	}
	if len(input.SpecificationName) > 200 {
		return input, errors.New("specification filename is too long")
	}
	input.Reviewed = false
	return input, nil
}

func safeText(s string) bool {
	if !utf8.ValidString(s) {
		return false
	}
	for _, r := range s {
		if unicode.IsControl(r) && r != '\n' && r != '\r' && r != '\t' {
			return false
		}
	}
	return true
}

type renderSection struct {
	Level      int      `json:"level"`
	Heading    string   `json:"heading"`
	Body       string   `json:"body"`
	Missing    []string `json:"missing"`
	Unnumbered bool     `json:"unnumbered"`
}

type renderData struct {
	Title         string            `json:"title"`
	Language      string            `json:"language"`
	Academic      bool              `json:"academic"`
	Metadata      map[string]string `json:"metadata"`
	DocumentCode  string            `json:"document_code"`
	ApprovalCode  string            `json:"approval_code"`
	Annotation    string            `json:"annotation"`
	Specification string            `json:"specification"`
	Sections      []renderSection   `json:"sections"`
	Labels        map[string]any    `json:"labels"`
}

func renderLabels(language string) map[string]any {
	if language == "ru" {
		return map[string]any{
			"agreed": "СОГЛАСОВАНО", "approved": "УТВЕРЖДЕНО", "approved_document": "УТВЕРЖДЕН", "approval_sheet": "ЛИСТ УТВЕРЖДЕНИЯ", "contributors": "Исполнители:", "group": "Группа:", "signature": "____________ / подпись /     Дата: ____________", "sheets": "Листов:", "annotation": "АННОТАЦИЯ", "contents": "СОДЕРЖАНИЕ", "changes": "ЛИСТ РЕГИСТРАЦИИ ИЗМЕНЕНИЙ", "specification": "Исходное техническое задание (данные пользователя)", "missing": "Недостающие сведения:", "review_warning": "Проект документа по проверенным пользователем разделам. Пробелы в исходных данных отмечены явно. Соответствие стандартам и фактическая корректность требуют проверки.", "inventory": []string{"Инв. № подл.", "Подп. и дата", "Взам. инв. №", "Инв. № дубл.", "Подп. и дата"}, "change_columns": []string{"Изм.", "Измененные листы", "Замененные листы", "Новые листы", "Аннулированные листы", "№ документа", "Подпись / дата"}}
	}
	return map[string]any{
		"agreed": "AGREED", "approved": "APPROVED", "approved_document": "APPROVED BY", "approval_sheet": "APPROVAL SHEET", "contributors": "Contributors:", "group": "Group:", "signature": "____________ / signature /     Date: ____________", "sheets": "Sheets:", "annotation": "ANNOTATION", "contents": "CONTENTS", "changes": "CHANGE REGISTRATION SHEET", "specification": "Original technical specification (user supplied)", "missing": "Missing information:", "review_warning": "Draft document from user-reviewed sections. Missing source material is marked explicitly. Standards compliance and factual accuracy require review.", "inventory": []string{"Original inventory no.", "Signature and date", "Replacement inventory no.", "Duplicate inventory no.", "Signature and date"}, "change_columns": []string{"Change", "Changed sheets", "Replaced sheets", "New sheets", "Cancelled sheets", "Document no.", "Signature / date"}}
}

// Render builds a private root containing trusted code and JSON data only.
func Render(ctx context.Context, destination string, job domain.Job, input domain.DocumentInput, sections []domain.Section, evidence ...domain.Evidence) ([]domain.Artifact, error) {
	p, err := templates.Load(job.ProfileID)
	if err != nil {
		return nil, err
	}
	if _, err = ValidateDocumentInput(input, job.ProfileID, job.Language); err != nil {
		return nil, err
	}
	if err = os.MkdirAll(destination, 0700); err != nil {
		return nil, err
	}
	root, err := os.MkdirTemp(destination, "render-")
	if err != nil {
		return nil, err
	}
	success := false
	defer func() {
		if !success {
			_ = os.RemoveAll(root)
		}
	}()
	missing := func(key string) string {
		if job.Language == "ru" {
			labels := map[string]string{"project": "название проекта", "organization": "организация", "faculty": "факультет", "program": "образовательная программа", "author": "автор или участники команды", "group": "учебная группа", "supervisor": "руководитель", "supervisor_position": "должность руководителя", "approver": "утверждающее лицо", "approver_position": "должность утверждающего лица", "year": "год", "code": "код проекта", "technical specification": "техническое задание"}
			return "[Не предоставлено: " + labels[key] + "]"
		}
		return "[Not provided: " + key + "]"
	}
	meta := map[string]string{}
	for _, key := range MetadataFields {
		meta[key] = input.Metadata[key]
		if meta[key] == "" {
			meta[key] = missing(key)
		}
	}
	data := renderData{Title: p.Title(job.Language), Language: job.Language, Academic: p.Academic, Metadata: meta, Labels: renderLabels(job.Language), Specification: input.Specification}
	if data.Specification == "" {
		data.Specification = missing("technical specification")
	}
	data.DocumentCode = "RU.17701729." + meta["code"] + " " + p.Code + " 01-1"
	if !p.Academic {
		data.DocumentCode = "Technical overview / " + p.Version
		if job.Language == "ru" {
			data.DocumentCode = "Технический обзор / " + p.Version
		}
	}
	data.ApprovalCode = data.DocumentCode + "-ЛУ"
	if job.Language == "en" {
		data.ApprovalCode = data.DocumentCode + "-AS"
	}
	data.Annotation = "This document describes the project using the section order of the selected profile. Source-linked drafts and user-supplied requirements are retained separately. No test execution or standards certification is implied."
	if job.Language == "ru" {
		data.Annotation = "Документ описывает проект в порядке разделов выбранного профиля. Черновики со ссылками на источники и требования пользователя сохранены раздельно. Выполнение испытаний и сертификация соответствия стандартам не подразумеваются."
	}
	indexed := map[string]domain.Section{}
	for _, s := range sections {
		indexed[s.ID] = s
	}
	for _, h := range p.Headings {
		item := renderSection{Level: h.Level, Heading: h.Title(job.Language), Missing: []string{}, Unnumbered: strings.Contains(strings.ToLower(h.EN), "appendix") || h.EN == "References"}
		if h.Draft {
			s, ok := indexed[h.ID]
			if !ok || s.Heading != h.Title(job.Language) {
				return nil, fmt.Errorf("reviewed section %s is missing or uses an obsolete plan", h.ID)
			}
			if !safeText(s.Body) {
				return nil, errors.New("section contains invalid text")
			}
			item.Body = s.Body
			item.Missing = s.Missing
		}
		data.Sections = append(data.Sections, item)
	}
	if err = fs.WalkDir(templates.Assets, "v1", func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() {
			return nil
		}
		relative := strings.TrimPrefix(path, "v1/")
		target := filepath.Join(root, relative)
		if err := os.MkdirAll(filepath.Dir(target), 0700); err != nil {
			return err
		}
		raw, err := templates.Assets.ReadFile(path)
		if err != nil {
			return err
		}
		return os.WriteFile(target, raw, 0600)
	}); err != nil {
		return nil, err
	}
	raw, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return nil, err
	}
	if err = os.WriteFile(filepath.Join(root, "document.json"), raw, 0600); err != nil {
		return nil, err
	}
	if len(evidence) > 0 {
		raw, err := json.MarshalIndent(evidence[0], "", "  ")
		if err != nil {
			return nil, err
		}
		if err = os.WriteFile(filepath.Join(root, "evidence.json"), raw, 0600); err != nil {
			return nil, err
		}
	}
	command := "typst compile --root . --font-path fonts --ignore-system-fonts --ignore-embedded-fonts --creation-timestamp 0 main.typ document.pdf"
	readme := fmt.Sprintf("# Editable document bundle\n\nRequires Typst %s. Edit document.json (text is always literal), or the trusted main.typ template, then run:\n\n%s\n\nFonts and their license are bundled. Profile version: %s. Only trusted local templates are compiled by the application. Do not compile arbitrary uploaded Typst.\n", templates.TypstVersion, command, p.Version)
	if err = os.WriteFile(filepath.Join(root, "README.md"), []byte(readme), 0600); err != nil {
		return nil, err
	}
	if err = CompileTypst(ctx, root); err != nil {
		return nil, err
	}
	bundle := filepath.Join(root, "source.zip")
	if err = zipBundle(root, bundle); err != nil {
		return nil, err
	}
	artifacts := []domain.Artifact{}
	for _, item := range []struct{ kind, name string }{{"pdf", "document.pdf"}, {"typst", "source.zip"}} {
		path := filepath.Join(root, item.name)
		raw, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}
		digest := sha256.Sum256(raw)
		artifacts = append(artifacts, domain.Artifact{Kind: item.kind, Path: path, SHA256: hex.EncodeToString(digest[:])})
	}
	success = true
	return artifacts, nil
}

func CompileTypst(parent context.Context, root string) error {
	binary, err := exec.LookPath("typst")
	if err != nil {
		return errors.New("Typst 0.15.1 is required on the worker")
	}
	ctx, cancel := context.WithTimeout(parent, 45*time.Second)
	defer cancel()
	version, err := exec.CommandContext(ctx, binary, "--version").Output()
	if err != nil || !strings.HasPrefix(string(version), "typst "+templates.TypstVersion+" ") {
		return errors.New("worker Typst version does not match pinned 0.15.1")
	}
	args := []string{"compile", "--root", root, "--font-path", filepath.Join(root, "fonts"), "--ignore-system-fonts", "--ignore-embedded-fonts", "--creation-timestamp", "0", "--package-path", filepath.Join(root, "packages"), "--package-cache-path", filepath.Join(root, "packages"), "main.typ", "document.pdf"}
	command := exec.CommandContext(ctx, binary, args...)
	command.Dir = root
	command.Env = []string{"HOME=" + root, "XDG_CACHE_HOME=" + root, "SOURCE_DATE_EPOCH=0"}
	var output bytes.Buffer
	command.Stdout = &output
	command.Stderr = &output
	if err = command.Run(); err != nil {
		return fmt.Errorf("compile trusted Typst template: %w: %.2000s", err, output.String())
	}
	return nil
}

func zipBundle(root, destination string) error {
	paths := []string{}
	if err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !entry.IsDir() && path != destination && !strings.HasSuffix(path, "document.pdf") {
			paths = append(paths, path)
		}
		return nil
	}); err != nil {
		return err
	}
	sort.Strings(paths)
	file, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer file.Close()
	writer := zip.NewWriter(file)
	for _, path := range paths {
		relative, _ := filepath.Rel(root, path)
		header := &zip.FileHeader{Name: filepath.ToSlash(relative), Method: zip.Deflate}
		header.SetModTime(time.Unix(0, 0))
		entry, err := writer.CreateHeader(header)
		if err != nil {
			return err
		}
		raw, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		if _, err = entry.Write(raw); err != nil {
			return err
		}
	}
	return writer.Close()
}

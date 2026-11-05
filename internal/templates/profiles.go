package templates

import (
	"embed"
	"encoding/json"
	"fmt"
)

//go:embed v1
var Assets embed.FS

const TypstVersion = "0.15.1"

type Heading struct {
	Chapter  string `json:"chapter,omitempty"`
	Guidance string `json:"guidance,omitempty"`
	ID       string `json:"id"`
	Level    int    `json:"level"`
	RU       string `json:"ru"`
	EN       string `json:"en"`
	CA       string `json:"ca,omitempty"`
	Draft    bool   `json:"draft"`
}

func (h Heading) Title(language string) string {
	if language == "ru" {
		return h.RU
	}
	if language == "ca" {
		return h.CA
	}
	return h.EN
}

type Profile struct {
	Guidance string    `json:"guidance,omitempty"`
	NoticeRU string    `json:"notice_ru,omitempty"`
	NoticeEN string    `json:"notice_en,omitempty"`
	NoticeCA string    `json:"notice_ca,omitempty"`
	ID       string    `json:"id"`
	Version  string    `json:"version"`
	Academic bool      `json:"academic"`
	RU       string    `json:"ru"`
	EN       string    `json:"en"`
	CA       string    `json:"ca,omitempty"`
	Code     string    `json:"code"`
	Headings []Heading `json:"headings"`
}

func (p Profile) Title(language string) string {
	if language == "ru" {
		return p.RU
	}
	if language == "ca" {
		return p.CA
	}
	return p.EN
}
func Load(id string) (Profile, error) {
	var profile Profile
	switch id {
	case "tz", "pz", "ro", "pmi", "pmi-team", "technical-overview", "srs", "arc42", "client-handover":
	default:
		return profile, fmt.Errorf("unknown document profile")
	}
	raw, err := Assets.ReadFile("v1/" + id + ".json")
	if err != nil {
		return profile, err
	}
	err = json.Unmarshal(raw, &profile)
	return profile, err
}

func (p Profile) Notice(language string) string {
	if language == "ru" {
		return p.NoticeRU
	}
	if language == "ca" {
		return p.NoticeCA
	}
	return p.NoticeEN
}

// SupportsLanguage requires complete fixed text before exposing Catalan support.
func (p Profile) SupportsLanguage(language string) bool {
	if language == "ru" || language == "en" {
		return true
	}
	if language != "ca" || p.Academic || p.CA == "" {
		return false
	}
	if (p.NoticeEN != "" || p.NoticeRU != "") && p.NoticeCA == "" {
		return false
	}
	for _, heading := range p.Headings {
		if heading.CA == "" {
			return false
		}
	}
	return true
}

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
	Draft    bool   `json:"draft"`
}

func (h Heading) Title(language string) string {
	if language == "ru" {
		return h.RU
	}
	return h.EN
}

type Profile struct {
	Guidance string    `json:"guidance,omitempty"`
	NoticeRU string    `json:"notice_ru,omitempty"`
	NoticeEN string    `json:"notice_en,omitempty"`
	ID       string    `json:"id"`
	Version  string    `json:"version"`
	Academic bool      `json:"academic"`
	RU       string    `json:"ru"`
	EN       string    `json:"en"`
	Code     string    `json:"code"`
	Headings []Heading `json:"headings"`
}

func (p Profile) Title(language string) string {
	if language == "ru" {
		return p.RU
	}
	return p.EN
}
func Load(id string) (Profile, error) {
	var profile Profile
	switch id {
	case "tz", "pz", "ro", "pmi", "pmi-team", "technical-overview", "srs", "arc42":
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
	return p.NoticeEN
}

package documents

import (
	"encoding/json"
	"path/filepath"
	"sort"
	"strings"
)

type manifestFact struct {
	Kind  string
	Value string
	Line  int
}

func manifestFacts(path string, data []byte) ([]manifestFact, bool) {
	if filepath.Base(path) == "package.json" {
		var manifest map[string]json.RawMessage
		if json.Unmarshal(data, &manifest) != nil {
			return nil, false
		}
		facts := []manifestFact{}
		for _, key := range []string{"name", "version", "description", "type", "main"} {
			var value string
			if json.Unmarshal(manifest[key], &value) == nil {
				facts = append(facts, manifestFact{"package_" + key, value, 0})
			}
		}
		for _, key := range []string{"scripts", "dependencies", "devDependencies", "peerDependencies", "engines"} {
			var values map[string]string
			if json.Unmarshal(manifest[key], &values) != nil {
				continue
			}
			names := make([]string, 0, len(values))
			for name := range values {
				names = append(names, name)
			}
			sort.Strings(names)
			for _, name := range names {
				facts = append(facts, manifestFact{"package_" + key, name + " = " + values[name], 0})
			}
		}

		return facts, true
	}
	if filepath.Base(path) != "go.mod" {
		return nil, false
	}
	facts := []manifestFact{}
	inRequire := false
	for index, line := range strings.Split(string(data), "\n") {
		value := strings.TrimSpace(strings.SplitN(line, "//", 2)[0])
		if value == "" {
			continue
		}
		if value == "require (" {
			inRequire = true
			continue
		}
		if value == ")" {
			inRequire = false
			continue
		}
		fields := strings.Fields(value)
		if len(fields) < 2 {
			continue
		}
		kind := ""
		switch fields[0] {
		case "module":
			kind = "go_module"
		case "go":
			kind = "go_version"
		case "toolchain":
			kind = "go_toolchain"
		case "require":
			kind = "go_dependency"
		}
		if inRequire {
			kind = "go_dependency"
		} else if kind != "" {
			value = strings.Join(fields[1:], " ")
		}
		if kind != "" {
			facts = append(facts, manifestFact{kind, value, index + 1})
		}
	}

	return facts, true
}

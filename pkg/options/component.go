package options

import (
	"errors"
	"fmt"
	"strings"
)

// ParseComponentOption parses --component options and interprets them as a map of components.
// key is component name, value is componentVersion
//
// eg options:
// --component A/v1.2.3 --component B/v0.0.2
//
// result:
//
// {
//   "A": "v1.2.3"
//   "B": "v0.0.2"
// }
func ParseComponentOption(components []string) (map[string]string, error) {
	out := make(map[string]string)

	for _, componentOpt := range components {
		s := strings.Split(componentOpt, "/")
		switch l := len(s); l {
		case 1:
			componentName := s[0]
			out[componentName] = ""
		case 2:
			componentName, componentVersion := s[0], s[1]
			out[componentName] = componentVersion
		case 3:
			componentName, _, componentGitVersion := s[0], s[1], s[2]
			out[componentName] = componentGitVersion
		default:
			msg := fmt.Sprintf("wrong --component option value, %s", componentOpt)
			return nil, errors.New(msg)
		}
	}

	return out, nil
}

// ParseComponentOption parses --components options and interprets them as a list of components.
// The value of --components option can be comma joined component names
// eg options:
// --components A,B --components C
//
// result:
// [A, B, C]
func ParseComponentsOption(components []string) []string {
	out := []string{}

	for _, componentsOpt := range components {
		s := strings.Split(componentsOpt, ",")
		out = append(out, s...)
	}

	return out
}

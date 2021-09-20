package options

import (
	"errors"
	"fmt"
	"strings"

	"github.com/bougou/sail/pkg/ansible"
)

// ParseHostsOptions parses --hosts options and interprets them as a map,
// key is component name, value is a list of ActionHosts.
// The returned map can be used to update cmdb inventory.
//
// eg options:
//    --hosts A,B/10.0.0.1,10.0.0.2 --hosts +C/10.0.0.3,10.0.0.4 --hosts -C,D,E/10.0.0.4
//
// result:
//
//    {
//      "A": [
//        { "Aciton": "update", "Hosts": ["10.0.0.1", "10.0.0.2"] },
//      ],
//      "B": [
//        { "Aciton": "update", "Hosts": ["10.0.0.1", "10.0.0.2"] },
//      ],
//      "C": [
//        { "Aciton": "add", "Hosts": ["10.0.0.3", "10.0.0.4"] },
//        { "Aciton": "remove", "Hosts": ["10.0.0.4"] },
//      ],
//      "D": [
//        { "Aciton": "remove", "Hosts": ["10.0.0.4"] },
//      ],
//      "E": [
//        { "Aciton": "remove", "Hosts": ["10.0.0.4"] },
//      ],
//    }
func ParseHostsOptions(hostsOptions []string) (map[string][]ansible.ActionHosts, error) {
	out := make(map[string][]ansible.ActionHosts)

	for _, hostsOpt := range hostsOptions {
		ah := ansible.ActionHosts{}

		var action string
		if strings.HasPrefix(hostsOpt, "-") {
			action = "remove"
			hostsOpt = strings.Replace(hostsOpt, "-", "", 1)
		} else if strings.HasPrefix(hostsOpt, "+") {
			action = "add"
			hostsOpt = strings.Replace(hostsOpt, "+", "", 1)
		} else {
			action = "update"
		}
		ah.Action = action

		s := strings.Split(hostsOpt, "/")
		switch l := len(s); l {
		case 1:
			ah.Hosts = strings.Split(s[0], ",")
			componentName := "_cluster"
			if _, exists := out[componentName]; !exists {
				out[componentName] = make([]ansible.ActionHosts, 0)
			}
			out[componentName] = append(out[componentName], ah)
		case 2:
			componentNames, hostsStr := s[0], s[1]
			for _, componentName := range strings.Split(componentNames, ",") {
				ah.Hosts = strings.Split(hostsStr, ",")
				if _, exists := out[componentName]; !exists {
					out[componentName] = make([]ansible.ActionHosts, 0)
				}
				out[componentName] = append(out[componentName], ah)
			}
		default:
			msg := fmt.Sprintf("wrong --hosts option value, %s", hostsOpt)
			return nil, errors.New(msg)
		}
	}

	return out, nil
}

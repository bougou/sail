package options

import (
	"errors"
	"fmt"
	"strings"

	"github.com/bougou/sail/pkg/models/product"
	"github.com/bougou/sail/pkg/models/target"
)

// ParseComponentsOption parses --component options and interprets them as map of components.
//
// eg options:
//   --component A                     # only specify componentName
//   --component B/v0.0.1              # specify componentName/componentVersion
//   --component C/v0.0.2,D/v0.0.3     # sepcify multiple components with comma separated
//   --component E/v1.2.3/831212f      # specify componentName/componentVersion/componentLongVersion
// result:
//
//   {
//     "A": "",
//     "B": "v0.0.1",
//     "C": "v0.0.2",
//     "D": "v0.0.3",
//     "E": "831212f"
//   }
func ParseComponentsOption(componentsOptions []string) (map[string]string, error) {
	out := make(map[string]string)

	for _, componentsOption := range componentsOptions {
		componentOpts := strings.Split(componentsOption, ",")

		for _, componentOpt := range componentOpts {
			s := strings.Split(componentOpt, "/")
			switch l := len(s); l {
			case 1:
				componentName := s[0]
				out[componentName] = ""
			case 2:
				componentName, componentVersion := s[0], s[1]
				out[componentName] = componentVersion
			case 3:
				componentName, _, componentLongVersion := s[0], s[1], s[2]
				out[componentName] = componentLongVersion
			default:
				msg := fmt.Sprintf("wrong --component option value, %s", componentOpt)
				return nil, errors.New(msg)
			}
		}

	}

	return out, nil
}

func ParseChoosedComponents(zone *target.Zone, components []string, ansible bool, helm bool) (map[string]string, map[string]string, error) {
	m, err := ParseComponentsOption(components)
	if err != nil {
		return nil, nil, fmt.Errorf("parse choosed components failed, err: %s", err)
	}

	for componentName, componentVersion := range m {
		if componentVersion != "" {
			zone.SetComponentVersion(componentName, componentVersion)
		}
	}

	if ansible {
		serverComponents := zone.Product.ComponentListWithFitlerOptionsOr(product.FilterOptionFormServer)
		for _, serverComponent := range serverComponents {
			if _, exists := m[serverComponent]; !exists {
				m[serverComponent] = ""
			}
		}
	}
	if helm {
		podComponents := zone.Product.ComponentListWithFitlerOptionsOr(product.FilterOptionFormPod)
		for _, podComponent := range podComponents {
			if _, exists := m[podComponent]; !exists {
				m[podComponent] = ""
			}
		}
	}

	serverComponents := make(map[string]string)
	podComponents := make(map[string]string)

	for k, v := range m {
		if zone.Product.HasComponent(k) {
			component, ok := zone.Product.Components[k]
			if !ok {
				return nil, nil, fmt.Errorf("not found component (%s) in zone", k)
			}
			switch component.Form {
			case product.ComponentFormServer:
				serverComponents[k] = v
			case product.ComponentFormPod:
				podComponents[k] = v
			default:
				serverComponents[k] = v
			}
		}
	}

	return serverComponents, podComponents, nil
}

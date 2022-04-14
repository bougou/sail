package ansible

import (
	"fmt"
	"strings"

	"github.com/kr/pretty"
	"gopkg.in/yaml.v3"
)

type Task struct {
	Name        string        `yaml:"name"`
	ChangedWhen bool          `yaml:"changed_when,omitempty"`
	FailedWhen  bool          `yaml:"failed_when,omitempty"`
	WithItems   *ListOrString `yaml:"with_items,omitempty"`

	Module `json:",inline" yaml:",inline"`
}

type Module map[string]ModuleArgs
type ModuleArgs map[string]interface{}

type ListOrString struct {
	list []interface{}
	s    string
}

func (l *ListOrString) MarshalYAML() (interface{}, error) {
	fmt.Println("call marshal")
	fmt.Println(len(l.list))
	if len(l.list) != 0 {
		ll := make([]string, len(l.list))
		for i, v := range l.list {

			ss, err := yaml.Marshal(v)
			if err != nil {
				return nil, err
			}
			fmt.Println("ss:", string(ss))
			ll[i] = strings.TrimSuffix(string(ss), "\n")
		}

		pretty.Println(ll)
		return yaml.Marshal(ll)

	}
	return []byte(l.s), nil
}

func (l *ListOrString) UnmarshalYAML(value *yaml.Node) error {
	pretty.Println(value)

	fmt.Println("call unmarshal")
	var list []interface{}
	var s string

	var ok bool

	if !ok {
		if err := value.Decode(&list); err == nil {
			ok = true
			l.list = list
		}
	}

	if !ok {
		if err := value.Decode(&s); err != nil {
			ok = true
			l.s = s
		}
	}

	if !ok {
		return fmt.Errorf("unmarshal failed")
	}

	return nil
}

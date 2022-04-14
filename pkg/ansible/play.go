package ansible

import "gopkg.in/yaml.v3"

type Play struct {
	Name           string    `yaml:"name"`
	Hosts          yaml.Node `yaml:"hosts"`
	AnyErrorsFatal bool      `yaml:"any_errors_fatal"`
	GatherFacts    bool      `yaml:"gather_facts"`
	Become         bool      `yaml:"become"`
	Tasks          []Task    `yaml:"tasks,omitempty"`
	Roles          []Role    `yaml:"roles,omitempty"`
	Tags           []string  `yaml:"tags,omitempty"`
}

func NewPlay(name string, hostsstr string) *Play {
	return &Play{
		Name: name,
		Hosts: yaml.Node{
			Kind:  yaml.ScalarNode,
			Style: yaml.DoubleQuotedStyle,
			Value: hostsstr,
		},
		GatherFacts:    true,
		AnyErrorsFatal: true,
		Tasks:          make([]Task, 0),
		Roles:          make([]Role, 0),
		Tags:           make([]string, 0),
	}
}

func (p *Play) AddTasks(tasks ...Task) *Play {
	p.Tasks = append(p.Tasks, tasks...)
	return p
}

func (p *Play) AddRoles(roles ...Role) *Play {
	p.Roles = append(p.Roles, roles...)
	return p
}

func (p *Play) AddTags(tags ...string) *Play {
	p.Tags = append(p.Tags, tags...)
	return p
}

func (p *Play) SetGatherFacts(flag bool) *Play {
	p.GatherFacts = flag
	return p
}

func (p *Play) SetAnyErrorsFatal(flag bool) *Play {
	p.AnyErrorsFatal = flag
	return p
}

type Role struct {
	Role string   `yaml:"role,omitempty"`
	Tags []string `yaml:"tags,omitempty"`
}

func (r *Role) WithTags(tags ...string) *Role {
	r.Tags = append(r.Tags, tags...)
	return r
}

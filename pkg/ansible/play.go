package ansible

import "gopkg.in/yaml.v3"

type Play struct {
	Name           string    `yaml:"name"`
	Hosts          yaml.Node `yaml:"hosts"`
	AnyErrorsFatal bool      `yaml:"any_errors_fatal"`
	GatherFacts    bool      `yaml:"gather_facts"`
	Become         bool      `yaml:"become"`
	Roles          []Role    `yaml:"roles"`
	Tags           []string  `yaml:"tags"`
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
	}
}

func (p *Play) WithRoles(roles ...Role) *Play {
	p.Roles = append(p.Roles, roles...)
	return p
}

func (p *Play) WithTags(tags ...string) *Play {
	p.Tags = append(p.Tags, tags...)
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

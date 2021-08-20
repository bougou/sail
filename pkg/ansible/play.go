package ansible

import "gopkg.in/yaml.v3"

type Play struct {
	Name           string    `yaml:"name,omitempty"`
	Hosts          yaml.Node `yaml:"hosts,omitempty"`
	AnyErrorsFatal bool      `yaml:"any_errors_fatal,omitempty"`
	GatherFacts    bool      `yaml:"gather_facts,omitempty"`
	Become         bool      `yaml:"become,omitempty"`
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
		GatherFacts: false,
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

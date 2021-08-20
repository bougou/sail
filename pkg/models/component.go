package models

import (
	"errors"
	"fmt"
	"strings"

	"github.com/bougou/sail/pkg/ansible"
	"gopkg.in/yaml.v3"
)

// Component represents the configuration of a component.
type Component struct {
	Name    string `yaml:"-"`
	Version string `yaml:"version"`

	// Pkgs hold all files that are used to complete the deployment of the component.
	Pkgs []Pkg `yaml:"pkgs"`

	Enabled  bool `yaml:"enabled"`
	External bool `yaml:"external"`

	Services map[string]Service         `yaml:"services"`
	Computed map[string]ServiceComputed `yaml:"computed"`

	Requires     []Require `yaml:"requires"`
	Dependencies []string  `yaml:"dependencies"`

	Roles []string `yaml:"roles"`

	Vars map[string]interface{} `yaml:"vars"`
	Tags map[string]interface{} `yaml:"tags"`
}

// NewComponent returns a Component.
func NewComponent(name string) *Component {
	return &Component{
		Name: name,

		Pkgs:         make([]Pkg, 0),
		Services:     make(map[string]Service),
		Computed:     make(map[string]ServiceComputed),
		Requires:     make([]Require, 0),
		Dependencies: make([]string, 0),
		Roles:        make([]string, 0),
		Vars:         make(map[string]interface{}),
		Tags:         make(map[string]interface{}),
	}
}

func (c *Component) Merge(in *Component) {

}

// Pkg represents a package file.
// We can use the Pkg.File field to check whether those pkg files exists
// and use the Pkg.URL field to download the file.
type Pkg struct {
	File *string `yaml:"file"`
	URL  *string `yaml:"url"`
}

type Require struct {
	Component *string `yaml:"component,omitempty"`
	Service   *string `yaml:"service,omitempty"`
}

func (c *Component) DownloadPkg(dstDir string) error {
	// Todo
	return nil
}

func (c *Component) Check() error {
	errs := []error{}
	for _, s := range c.Services {
		if err := s.Check(); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) == 0 {
		return nil
	}

	errmsgs := []string{}
	for _, err := range errs {
		errmsgs = append(errmsgs, err.Error())
	}

	msg := fmt.Sprintf("Check component (%s) external (%v) failed, err: %s", c.Name, c.External, strings.Join(errmsgs, "; "))
	return errors.New(msg)
}

func (c *Component) Compute(cmdb *CMDB) error {
	// Todo
	for svcName, svc := range c.Services {
		svcComputed, err := svc.Compute(c.External, cmdb)
		if err != nil {
			msg := fmt.Sprintf("Compute service (%s) failed, err: %s", svcName, err)
			return errors.New(msg)
		}

		(c.Computed)[svcName] = *svcComputed
	}
	return nil
}

func (c *Component) GenAnsiblePlay() (*ansible.Play, error) {
	hostsPattern := fmt.Sprintf("{{ _ansiblepattern_%s | default('%s') }}", strings.ReplaceAll(c.Name, "-", "_"), c.Name)
	play := ansible.NewPlay(c.Name, hostsPattern)
	play.WithTags("hosts-" + c.Name)

	if len(c.Roles) == 0 {
		role := ansible.Role{
			Role: c.Name,
			Tags: []string{c.Name},
		}
		play.Roles = append(play.Roles, role)
	} else {
		for _, r := range c.Roles {
			role := ansible.Role{
				Role: r,
				Tags: []string{r},
			}
			play.Roles = append(play.Roles, role)
		}
	}

	return play, nil
}

func newComponentFromValue(componentName string, componentValue interface{}) (*Component, error) {
	b, err := yaml.Marshal(componentValue)
	if err != nil {
		msg := fmt.Sprintf("marshal failed, err: %s", err)
		return nil, errors.New(msg)
	}

	c := NewComponent(componentName)
	if err := yaml.Unmarshal(b, c); err != nil {
		msg := fmt.Sprintf("yaml unmarshal failed, err: %s", err)
		return nil, errors.New(msg)
	}

	if c.Enabled && c.External {
		msg := fmt.Sprintf("Warn: Enabled and External of component can not be both true, automatically set Enabled to false for component (%s)", componentName)
		fmt.Println(msg)
		c.Enabled = false
	}

	for svcName, s := range c.Services {
		b, err := yaml.Marshal(s)
		if err != nil {
			msg := fmt.Sprintf("marshal failed, err: %s", err)
			return nil, errors.New(msg)
		}

		outs := NewService(componentName, svcName)
		if err := yaml.Unmarshal(b, outs); err != nil {
			msg := fmt.Sprintf("yaml unmarshal failed, err: %s", err)
			return nil, errors.New(msg)
		}

		(c.Services)[svcName] = *outs
	}

	c.Name = componentName

	return c, nil
}

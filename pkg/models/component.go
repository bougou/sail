package models

import (
	"errors"
	"fmt"
	"strings"

	"github.com/bougou/sail/pkg/ansible"
	"gopkg.in/yaml.v3"
)

const (
	ComponentFormPod    = "pod"
	ComponentFormServer = "server"
)

// Component represents the configuration of a component.
type Component struct {
	Name    string `yaml:"-"`
	Version string `yaml:"version"`

	// If RoleName is empty, it will be set to same with Name.
	RoleName string `yaml:"roleName"`

	// Form represents the installation method of this component (组件的部署形态), valid values:
	//  * server
	//  * pod
	Form string `yaml:"form"`

	// Pkgs hold all files that are used to complete the deployment of the component.
	Pkgs []Pkg `yaml:"pkgs"`

	// Enabled represents wether this component is deployed in the specific environment.
	Enabled bool `yaml:"enabled"`

	// External represents wether this component is provided by external system like cloud, and thus no need to be deployed.
	External bool `yaml:"external"`

	// Services holds all exposed service
	Services map[string]Service `yaml:"services"`

	// Computed holds all auto computed service info
	Computed map[string]ServiceComputed `yaml:"computed"`

	// Requires represents the other components on which this component depends on.
	// If this component is activated (enabled:true or external:true), then all these required components also need to be activated.
	Requires []Require `yaml:"requires"`

	// The list value of `deps` represents the other services which depend on this service.
	// If the number of hosts of this service changed, it required that
	// those services who depend on it also need to be reconfigured or restarted.
	// For `service-scaleup` and `service-scaledown`, these dependencies will be used
	Dependencies []string `yaml:"dependencies"`

	// Children represents some children-level components of this component.
	// The children components have the following characteristics:
	// * They can declared to be activated or not enabled: true or false.
	// * They can be upgraded like normal components.
	// * They DO NOT have their own hosts group in ansible inventory, they share the hosts group of its parent component.
	Children []string `yaml:"children"`

	// Applied roles for this component.
	// The empty list will apply at least one role with the componen's RoleName.
	Roles []string `yaml:"roles"`

	Vars map[string]interface{} `yaml:"vars"`
	Tags map[string]interface{} `yaml:"tags"`
}

// NewComponent returns a Component.
func NewComponent(name string) *Component {
	return &Component{
		Name:     name,
		RoleName: name,

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

func (c *Component) GetRoleName() string {
	if c.RoleName != "" {
		return c.RoleName
	}
	return c.Name
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

	msg := fmt.Sprintf("check component (%s) external (%v) failed, err: %s", c.Name, c.External, strings.Join(errmsgs, "; "))
	return errors.New(msg)
}

func (c *Component) Compute(cmdb *CMDB) error {
	// Todo
	for svcName, svc := range c.Services {
		svcComputed, err := svc.Compute(c.External, cmdb)
		if err != nil {
			msg := fmt.Sprintf("compute service (%s) failed, err: %s", svcName, err)
			return errors.New(msg)
		}

		(c.Computed)[svcName] = *svcComputed
	}
	return nil
}

func (c *Component) GetRoles() []string {
	roles := []string{}
	if len(c.Roles) == 0 {
		roles = append(roles, c.RoleName)
	}
	for _, role := range c.Roles {
		if role == "." {
			role = c.RoleName
		}
		roles = append(roles, role)
	}
	return roles
}

func (c *Component) GenAnsiblePlay() (*ansible.Play, error) {
	hostsPattern := fmt.Sprintf("{{ _ansiblepattern_%s | default('%s') }}", strings.ReplaceAll(c.Name, "-", "_"), c.Name)
	play := ansible.NewPlay(c.Name, hostsPattern)
	play.WithTags("play-" + c.Name)

	for _, r := range c.GetRoles() {
		role := ansible.Role{
			Role: r,
			Tags: []string{r},
		}
		play.Roles = append(play.Roles, role)
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
		msg := fmt.Sprintf("warn: Enabled and External of component can not be both true, automatically set Enabled to false for component (%s)", componentName)
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

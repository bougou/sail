package product

import (
	"fmt"
	"strings"

	"github.com/bougou/sail/pkg/ansible"
	"github.com/bougou/sail/pkg/models/cmdb"
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

	// Group can be used to group by components.
	Group string `yaml:"group"`

	// If RoleName is empty, it will be set to component Name.
	RoleName string `yaml:"roleName"`

	// Form represents the installation method of this component, valid values:
	//  * server
	//  * pod
	// (组件的部署形态)
	Form string `yaml:"form"`

	// Pkgs hold all files that are used to complete the deployment of the component.
	Pkgs []Pkg `yaml:"pkgs"`

	// Enabled represents whether this component will be deployed in the specific environment.
	Enabled bool `yaml:"enabled"`

	// External represents whether this component is provided by external system like cloud, and thus no need to be deployed.
	External bool `yaml:"external"`

	// Services holds all exposed service of the component.
	// Each service is exposed by a specific port.
	Services map[string]Service `yaml:"services"`

	// Computed holds all auto computed service info.
	// This field SHOULD NEVER be edited or changed by operators.
	// The field is always automaticllay computed based on other fields of Component.
	Computed map[string]ServiceComputed `yaml:"computed"`

	// Requires represents the other components on which this component depends on.
	// If this component is activated (enabled:true or external:true), then all these required components also need to be activated.
	// 依赖的服务（其它组件提供的服务，不能依赖自身组件）
	// Todo, check cycle
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
	// The empty list will apply at least one role with the component's RoleName.
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

	return fmt.Errorf("check component (%s) external (%v) failed, err: %s", c.Name, c.External, strings.Join(errmsgs, "; "))
}

func (c *Component) Compute(cm *cmdb.CMDB) error {
	// Todo
	for svcName, svc := range c.Services {
		svcComputed, err := svc.Compute(c.External, cm)
		if err != nil {
			return fmt.Errorf("compute service (%s) failed, err: %s", svcName, err)
		}

		(c.Computed)[svcName] = *svcComputed
	}
	return nil
}

// GetRoles return the ansible roles which will be applied to this component.
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

// GenAnsiblePlay generatea a ansible play for this component.
func (c *Component) GenAnsiblePlay() (*ansible.Play, error) {
	// Note, we use the component name as the default ansible group name.
	hostsPattern := fmt.Sprintf("{{ _ansiblepattern_%s | default('%s') }}", strings.ReplaceAll(c.Name, "-", "_"), c.Name)
	play := ansible.NewPlay(c.Name, hostsPattern)
	play.AddTags("play-" + c.Name)
	if c.Group != "" {
		play.AddTags(c.Group)
	}

	for _, roleName := range c.GetRoles() {
		role := ansible.Role{
			Role: roleName,
			Tags: []string{roleName},
		}
		if roleName == c.RoleName && roleName != c.Name {
			role.Tags = append(role.Tags, c.Name)
		} else {
			role.Tags = append(role.Tags, roleName+"-"+c.Name)
		}

		play.Roles = append(play.Roles, role)
	}

	return play, nil
}

func newComponentFromValue(componentName string, componentValue interface{}) (*Component, error) {
	b, err := yaml.Marshal(componentValue)
	if err != nil {
		return nil, fmt.Errorf("marshal failed, err: %s", err)
	}

	c := NewComponent(componentName)
	if err := yaml.Unmarshal(b, c); err != nil {
		return nil, fmt.Errorf("yaml unmarshal failed, err: %s", err)
	}

	for svcName, s := range c.Services {
		b, err := yaml.Marshal(s)
		if err != nil {
			return nil, fmt.Errorf("marshal failed, err: %s", err)
		}

		outs := NewService(componentName, svcName)
		if err := yaml.Unmarshal(b, outs); err != nil {
			return nil, fmt.Errorf("yaml unmarshal failed, err: %s", err)
		}

		(c.Services)[svcName] = *outs
	}

	c.Name = componentName

	return c, nil
}

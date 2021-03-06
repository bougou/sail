package cmdb

import (
	"path"
	"strings"

	"github.com/bougou/sail/pkg/ansible"
	"github.com/mitchellh/go-homedir"
)

type CMDB struct {
	Inventory *ansible.Inventory  `yaml:"inventory"` // Ansible 格式的主机清单
	Platforms map[string]Platform `yaml:"platforms"` // 非主机部署形态, map key is component name or 'all'
}

type Platform struct {
	K8S *K8S `yaml:"k8s,omitempty"`
}

type K8S struct {
	KubeConfig  string `yaml:"kubeConfig"`
	KubeContext string `yaml:"kubeContext"`
	Namespace   string `yaml:"namespace"`
}

func ExpandTilde(pathstr string) string {
	if strings.HasPrefix(pathstr, "~") {
		home, err := homedir.Dir()
		if err != nil {
			return ""
		}
		s := strings.Replace(pathstr, "~", "", 1)
		return path.Join(home, s)
	}

	return pathstr
}

func NewCMDB() *CMDB {
	i := ansible.NewAnsibleInventory()
	i.FillAll()

	return &CMDB{
		Inventory: i,
		Platforms: make(map[string]Platform),
	}
}

func (c *CMDB) GetHostsForComponent(name string) []string {
	g, err := c.Inventory.GetGroup(name)
	if err != nil {
		return []string{}
	}

	return g.HostsList()
}

func (c *CMDB) Compute(componentName string, componentEnabled bool) error {
	if !componentEnabled {
		c.Inventory.RemoveGroup(componentName)
		return nil
	}

	if c.Inventory.HasGroup(componentName) {
		return nil
	}

	compHosts := c.determineHostsForComponent(componentName)
	group := ansible.NewGroup(componentName)
	group.AddHosts(compHosts...)
	return c.Inventory.AddGroup(group)
}

// Todo, recursively
func (c *CMDB) determineHostsForComponent(componentName string) []string {
	if c.Inventory.HasGroup("_cluster") {
		clusterGroup, _ := c.Inventory.GetGroup("_cluster")
		return clusterGroup.HostsList()
	}

	return []string{}
}

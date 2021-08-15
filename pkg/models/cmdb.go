package models

import "github.com/bougou/sail/pkg/ansible"

type CMDB struct {
	Inventory *ansible.Inventory     `yaml:"inventory"` // Ansible 格式的主机清单
	Platforms map[string]interface{} `yaml:"platforms"` // 非主机部署形态
}

func NewCMDB() *CMDB {
	i := ansible.NewAnsibleInventory()
	i.FillAll()

	return &CMDB{
		Inventory: i,
		Platforms: make(map[string]interface{}),
	}
}

func (c *CMDB) GetHostsForComponent(name string) []string {
	g, err := c.Inventory.GetGroup(name)
	if err != nil {
		return []string{}
	}

	return g.HostsList()
}

func (c *CMDB) Compute(components map[string]*Component) error {

	for compKey, compValue := range components {
		if !compValue.Enabled {
			c.Inventory.RemoveGroup(compKey)
			continue
		}

		if c.Inventory.HasGroup(compKey) {
			continue
		}

		compHosts := c.determineHostsForComponent(compKey)
		group := ansible.NewGroup(compKey)
		group.AddHosts(compHosts...)
		c.Inventory.AddGroup(group)
	}
	return nil
}

// Todo, recursively
func (c *CMDB) determineHostsForComponent(componentName string) []string {
	if c.Inventory.HasGroup("_cluster") {
		clusterGroup, _ := c.Inventory.GetGroup("_cluster")
		return clusterGroup.HostsList()
	}

	return []string{}
}

package ansible

const (
	MetaGroupName      string = "_meta"
	AllGroupName       string = "all"
	UngroupedGroupName string = "ungrouped"
)

type Group struct {
	name string `json:"-" yaml:"-"`

	// Hosts store hostvars for each host
	Hosts *map[string]map[string]interface{} `json:"hosts,omitempty" yaml:"hosts,omitempty"`

	// Vars store groupvars for this group
	Vars *map[string]interface{} `json:"vars,omitempty" yaml:"vars,omitempty"`

	Children *Inventory `json:"children,omitempty" yaml:"children,omitempty"`
}

// NewGroup create a noraml group with name
func NewGroup(name string) *Group {
	vars := make(map[string]interface{})
	hosts := make(map[string]map[string]interface{})

	return &Group{
		name: name,

		Hosts:    &hosts,
		Vars:     &(vars),
		Children: NewInventory(),
	}
}

func (m *Group) SetHostVar(host, varName string, varValue interface{}) {
	if _, ok := (*m.Hosts)[host]; !ok {
		(*m.Hosts)[host] = map[string]interface{}{}
	}
	(*m.Hosts)[host][varName] = varValue
}

func (m *Group) SetHostVars(host string, vars map[string]interface{}) {
	for varName, varValue := range vars {
		m.SetHostVar(host, varName, varValue)
	}
}

func (g *Group) Name() string {
	return g.name
}

func (g *Group) HasHost(host string) bool {
	_, ok := (*g.Hosts)[host]
	return ok
}

func (g *Group) AddHost(host string) {
	if !g.HasHost(host) {
		hostvars := make(map[string]interface{})
		(*g.Hosts)[host] = hostvars
	}
}

func (g *Group) AddHosts(hosts ...string) {
	for _, host := range hosts {
		g.AddHost(host)
	}
}

func (g *Group) RemoveHost(host string) {
	delete(*g.Hosts, host)
}

func (g *Group) RemoveHosts(hosts ...string) {
	for _, host := range hosts {
		g.RemoveHost(host)
	}
}

func (g *Group) HostsList() []string {
	out := []string{}
	for host := range *g.Hosts {
		out = append(out, host)
	}
	return out
}

func (g *Group) AddVar(key string, value interface{}) {
	(*g.Vars)[key] = value
}

func (g *Group) AddVars(vars map[string]interface{}) {
	for k, v := range vars {
		g.AddVar(k, v)
	}
}

func (g *Group) RemoveVar(key string) {
	delete(*g.Vars, key)
}

// SetChildren set inventory as the Children of group by overriding the Children field
func (g *Group) SetChildren(inventory *Inventory) {
	g.Children = inventory
}

// AddChildren set the groups in inventory to Children of group.
func (g *Group) AddChildren(inventory *Inventory) {
	for _, group := range inventory.GroupsMap {
		g.AddChildGroup(group)
	}
}

func (g *Group) AddChildGroup(group *Group) {
	g.Children.SetGroup(group)
}

// Todo, support golang generic
func unique(array []string) []string {
	keys := make(map[string]bool)
	out := []string{}
	for _, item := range array {
		if _, ok := keys[item]; !ok {
			keys[item] = true
			out = append(out, item)
		}
	}
	return out
}

func (g *Group) AddDefaultVars() {
	g.AddVars(map[string]interface{}{
		"ansible_port": 22,
		"ansible_user": "root",
		// "ssh_mode":         "password",
		// "ansible_password": "",
	})
}

func (g *Group) SetDefaultSSHUser(user string) {
	g.AddVar("ansible_user", user)
}

func (g *Group) SetDefaultSSHPort(port int) {
	g.AddVar("ansible_port", port)
}

func (g *Group) SetDefaultSSHPassword(password string) {
	g.AddVar("ansible_password", password)
}

func (g *Group) SetDefaultSSHMode(mode string) {
	g.AddVar("ssh_mode", mode)
}

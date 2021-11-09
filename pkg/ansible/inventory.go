package ansible

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"strings"
	"sync"

	"gopkg.in/yaml.v3"
)

type GroupsMap map[string]*Group

type Inventory struct {
	GroupsMap `json:",inline" yaml:",inline"`
	lock      sync.Mutex `json:"-" yaml:"-"`
}

func NewInventory() *Inventory {
	return &Inventory{
		GroupsMap: make(map[string]*Group),
		lock:      sync.Mutex{},
	}
}

// NewAnsibleInventory create a top level inventory.
// You should call NewInventory to create other level inventory.
func NewAnsibleInventory() *Inventory {
	i := NewInventory()
	i.AddAllGroup()

	return i
}

// AddAllGroup adds a "all" group to the inventory.
// Only the top level inventory should call this method.
func (i *Inventory) AddAllGroup() {
	g := NewGroup(AllGroupName)

	// just ingore error check for existed "all" group
	_ = i.AddGroup(g)
}

// MarshalJSON provides custom method when json.Marshal(i).
// 'inline' is not an official option of JSON struct tags.
func (i *Inventory) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.GroupsMap)
}

func (i *Inventory) UnmarshalJSON(data []byte) error {
	var m GroupsMap = make(map[string]*Group)
	decoder := json.NewDecoder(bytes.NewBuffer(data))
	// decoder.UseNumber()
	if err := decoder.Decode(&m); err != nil {
		return err
	}
	i.GroupsMap = m
	return nil
}

func (i *Inventory) UnmarshalYAML(value *yaml.Node) error {
	var m GroupsMap = make(map[string]*Group)
	if err := value.Decode(&m); err != nil {
		return err
	}
	i.GroupsMap = m
	return nil
}

func (i *Inventory) HasGroup(groupName string) bool {
	_, exist := i.GroupsMap[groupName]
	return exist
}

func (i *Inventory) HasMetaGroup() bool {
	return i.HasGroup(MetaGroupName)
}

func (i *Inventory) HasAllGroup() bool {
	return i.HasGroup(AllGroupName)
}

func (i *Inventory) GetGroup(groupName string) (*Group, error) {
	if i.HasGroup(groupName) {
		return i.GroupsMap[groupName], nil
	}

	return nil, errors.New("group not exist")
}

// SetGroup set group to inventory, override if group already exist
func (i *Inventory) SetGroup(group *Group) {
	i.lock.Lock()
	defer i.lock.Unlock()

	name := group.Name()
	i.GroupsMap[name] = group
}

// SetGroup set group to inventory, override if group already exist
func (i *Inventory) RemoveGroup(groupName string) {
	i.lock.Lock()
	defer i.lock.Unlock()

	delete(i.GroupsMap, groupName)
}

// AddGroup add group to inventory, return error if group already exist
func (i *Inventory) AddGroup(group *Group) error {
	if i.HasGroup(group.Name()) {
		return errors.New("group alread exist")
	}

	i.lock.Lock()
	defer i.lock.Unlock()

	name := group.Name()
	i.GroupsMap[name] = group
	return nil
}

func (i *Inventory) Merge(new *Inventory) error {
	if new == nil {
		return nil
	}

	for _, group := range new.GroupsMap {
		i.SetGroup(group)
	}

	return nil
}

func (i *Inventory) FilterOutIP(ipnet *net.IPNet) {
	if ipnet == nil {
		return
	}

	for _, group := range i.GroupsMap {
		for hostName := range *group.Hosts {
			ip := net.ParseIP(hostName)
			if !ipnet.Contains(ip) {
				group.RemoveHost(hostName)
			}
		}

		if group.Children != nil {
			group.Children.FilterOutIP(ipnet)
		}
	}
}

type InventoryFinderFunc func() (*Inventory, error)

func (i *Inventory) MergeWithFuncs(inventoryFinderFuncs ...InventoryFinderFunc) error {
	var errs []error
	var msgs []string

	for _, f := range inventoryFinderFuncs {
		inventory, err := f()
		if err != nil {
			errs = append(errs, err)
			msg := fmt.Sprintf("call inventoryFinderFunc (%v) failed, err: %s", f, err)
			msgs = append(msgs, msg)
			continue
		}

		if err := i.Merge(inventory); err != nil {
			errs = append(errs, err)
			msg := fmt.Sprintf("merge inventoryFinderFunc (%v) got inventory failed, err: %s", f, err)
			msgs = append(msgs, msg)
			continue
		}
	}

	if len(errs) != 0 {
		return errors.New(strings.Join(msgs, "\n"))
	}

	return nil
}

func (i *Inventory) GetAllHosts() []string {
	out := []string{}

	for groupName, group := range i.GroupsMap {
		if groupName == AllGroupName || groupName == MetaGroupName {
			continue
		}

		for host := range *group.Hosts {
			out = append(out, host)
		}

		if group.Children != nil {
			childrenHosts := (*group.Children).GetAllHosts()
			out = append(out, childrenHosts...)
		}
	}

	return out
}

// FillAll adds default vars to the 'all' group of the inventory.
// Only the top level inventory should have 'all' group, so should call this method.
func (i *Inventory) FillAll() {
	if !i.HasAllGroup() {
		return
	}

	a, _ := i.GetGroup(AllGroupName)
	allHosts := i.GetAllHosts()

	for _, host := range allHosts {
		a.AddHost(host)
	}

	a.AddDefaultVars()

}

func (i *Inventory) SetDefaultSSHUser(user string) {
	if !i.HasAllGroup() {
		return
	}
	a, _ := i.GetGroup(AllGroupName)
	a.AddVar("ansible_user", user)
}

func (i *Inventory) SetDefaultSSHPort(port int) {
	if !i.HasAllGroup() {
		return
	}
	a, _ := i.GetGroup(AllGroupName)
	a.AddVar("ansible_port", port)
}

func (i *Inventory) SetDefaultSSHPassword(password string) {
	if !i.HasAllGroup() {
		return
	}
	a, _ := i.GetGroup(AllGroupName)
	a.AddVar("ansible_password", password)
}

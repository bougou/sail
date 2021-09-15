package models

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strings"

	"log"

	"github.com/bougou/gopkg/common"
	"github.com/bougou/sail/pkg/ansible"
	"gopkg.in/yaml.v3"
)

const (
	SailMetaVarProduct  = "_sail_product"
	SailMetaVarHelmMode = "_sail_helm_mode"

	SailHelmModeComponent = "component"
	SailHelmModeProduct   = "product"
)

type ZoneMeta struct {
	// Product Name
	SailProduct  string `json:"_sail_product" yaml:"_sail_product"`     // tag value must equal to SailMetaVarProduct
	SailHelmMode string `json:"_sail_helm_mode" yaml:"_sail_helm_mode"` // tag value must equal to SailMetaVarHelmMode
}

type Zone struct {
	*ZoneMeta

	TargetName   string
	ZoneName     string
	TargetDir    string
	ZoneDir      string
	VarsFile     string
	HostsFile    string
	ComputedFile string
	ResourcesDir string

	HelmDir string

	Product  *Product
	CMDB     *CMDB
	Computed map[string]interface{}

	ansibleCfgFile string

	sailOption *SailOption
}

func NewZone(sailOption *SailOption, targetName string, zoneName string) *Zone {
	zone := &Zone{
		TargetName: targetName,
		ZoneName:   zoneName,

		TargetDir: path.Join(sailOption.TargetsDir, targetName),
		ZoneDir:   path.Join(sailOption.TargetsDir, targetName, zoneName),

		HostsFile:    path.Join(sailOption.TargetsDir, targetName, zoneName, "hosts.yaml"),
		VarsFile:     path.Join(sailOption.TargetsDir, targetName, zoneName, "vars.yaml"),
		ComputedFile: path.Join(sailOption.TargetsDir, targetName, zoneName, "_computed.yaml"),

		ResourcesDir: path.Join(sailOption.TargetsDir, targetName, zoneName, "resources"),

		HelmDir: path.Join(sailOption.TargetsDir, targetName, zoneName, "helm"),

		CMDB:     NewCMDB(),
		Computed: make(map[string]interface{}),

		ansibleCfgFile: path.Join(sailOption.ProductsDir, "ansible.cfg"),

		sailOption: sailOption,
	}

	return zone
}

// LoadNew fill vars to zone. The zone is treated as a newly created zone.
// So it will ONLY load default varibles from product.
// This method is ONLY called when `conf-create`.
func (zone *Zone) LoadNew() error {
	// for newly created zone, the zone.ProductName is set by conf-create
	if zone.SailProduct == "" {
		return errors.New("empty product name")
	}
	product := NewProduct(zone.SailProduct, zone.sailOption.ProductsDir)
	if err := product.Init(); err != nil {
		msg := fmt.Sprintf("init product failed, err: %s", err)
		return errors.New(msg)
	}

	zone.Product = product

	// fill zone meta vars
	zone.Product.Vars[SailMetaVarProduct] = zone.SailProduct
	zone.Product.Vars[SailMetaVarHelmMode] = zone.SailHelmMode

	return nil
}

// Load initialize the zone. The zone is supposed to be already exists.
// It will try to determine the product name from zone vars file.
func (zone *Zone) Load() error {
	zoneMeta, err := zone.ParseZoneMeta()
	if err != nil {
		return fmt.Errorf("parse zone meta failed, err: %s", err)
	}
	zone.ZoneMeta = zoneMeta

	if zone.SailProduct == "" {
		return errors.New("empty product name")
	}

	product := NewProduct(zone.SailProduct, zone.sailOption.ProductsDir)
	if err := product.Init(); err != nil {
		msg := fmt.Sprintf("init product failed, err: %s", err)
		return errors.New(msg)
	}

	if err := zone.LoadHosts(); err != nil {
		msg := fmt.Sprintf("load hosts failed, err: %s", err)
		return errors.New(msg)
	}

	if err := product.LoadZone(zone.VarsFile); err != nil {
		msg := fmt.Sprintf("load zone vars failed, err: %s", err)
		return errors.New(msg)
	}

	zone.Product = product

	if err := zone.PrepareHelm(); err != nil {
		msg := fmt.Sprintf("prepare helm failed, err: %s", err)
		return errors.New(msg)
	}

	return nil
}

func (zone *Zone) Compute() error {
	if err := zone.CMDB.Compute(zone.Product.Components); err != nil {
		msg := fmt.Sprintf("compute zone CMDB failed, err: %s", err)
		return errors.New(msg)
	}

	if err := zone.Product.Compute(zone.CMDB); err != nil {
		msg := fmt.Sprintf("compute zone product failed, err: %s", err)
		return errors.New(msg)
	}
	return nil
}

func (zone *Zone) ParseZoneMeta() (*ZoneMeta, error) {
	b, err := os.ReadFile(zone.VarsFile)
	if err != nil {
		msg := fmt.Sprintf("read zone vars file failed, err: %s", err)
		return nil, errors.New(msg)
	}

	m := &ZoneMeta{}
	if err := yaml.Unmarshal(b, &m); err != nil {
		msg := fmt.Sprintf("yaml unmarshal failed, err: %s", err)
		return nil, errors.New(msg)
	}

	if m.SailProduct == "" {
		msg := fmt.Sprintf("not found (%s) variable in %s, you have to fix that before continue", SailMetaVarProduct, zone.VarsFile)
		return nil, errors.New(msg)
	}

	return m, nil
}

func (zone *Zone) HandleCompatibity() {
	// Todo
	// Domain Specific Language (Declarative)
	// migrate.yaml
}

func (zone *Zone) Dump() error {
	if err := zone.Compute(); err != nil {
		msg := fmt.Sprintf("zone.Compute failed, err: %s", err)
		return errors.New(msg)
	}

	if err := os.MkdirAll(zone.ZoneDir, os.ModePerm); err != nil {
		msg := fmt.Sprintf("make zone dir failed, err: %s", err)
		return errors.New(msg)
	}

	zone.RenderVars()
	zone.RenderHosts()
	zone.RenderComputed()
	return nil
}

func (zone *Zone) RenderVars() {
	m := make(map[string]interface{})

	for k, v := range zone.Product.Vars {
		m[k] = v
	}

	for k, v := range zone.Product.Components {
		m[k] = v
	}

	b, err := common.Encode("yaml", m)
	if err != nil {
		fmt.Println("encode vars failed", err)
	}

	if err := os.WriteFile(zone.VarsFile, b, 0644); err != nil {
		fmt.Println("write vars file failed", err)
	}
}

func (zone *Zone) RenderHosts() {
	b, err := common.Encode("yaml", zone.CMDB.Inventory)
	if err != nil {
		fmt.Println("encode vars failed", err)
	}

	if err := os.WriteFile(zone.HostsFile, b, 0644); err != nil {
		fmt.Println("write hosts file failed", err)
	}
}

func (zone *Zone) RenderComputed() {
	b, err := common.Encode("yaml", zone.Computed)
	if err != nil {
		fmt.Println("encode vars failed", err)
	}

	if err := os.WriteFile(zone.ComputedFile, b, 0644); err != nil {
		fmt.Println("write hosts file failed", err)
	}
}

func (zone *Zone) PatchActionHostsMap(m map[string][]ActionHosts) error {
	for groupName, ahs := range m {
		if !zone.Product.HasComponent(groupName) && groupName != "_cluster" {
			msg := fmt.Sprintf("not supported component in this product, supported components: %s", zone.Product.ComponentList())
			return errors.New(msg)
		}

		for _, ah := range ahs {
			zone.PatchActionHosts(groupName, &ah)
		}
	}

	return nil
}

func (zone *Zone) PatchActionHosts(groupName string, hostsPatch *ActionHosts) {
	if zone.CMDB.Inventory.HasGroup(groupName) {
		group, _ := zone.CMDB.Inventory.GetGroup(groupName)
		PatchAnsibleGroup(group, hostsPatch)
	} else {
		if hostsPatch.Action == "delete" {
			return
		}

		group := ansible.NewGroup(groupName)
		for _, host := range hostsPatch.Hosts {
			group.AddHost(host)
			group.SetHostVars(host, map[string]interface{}{})
		}
		zone.CMDB.Inventory.AddGroup(group)
	}
}

func (zone *Zone) BuildInventory(hostsMap map[string][]string) error {

	for k, v := range hostsMap {
		if !zone.Product.HasComponent(k) {
			log.Printf("%s is not valid components, omit, valid components: %v\n", k, zone.Product.ComponentList())
			continue
		}

		group := ansible.NewGroup(k)
		for _, host := range v {
			group.AddHost(host)
			group.SetHostVars(host, map[string]interface{}{})
		}

		zone.CMDB.Inventory.AddGroup(group)
	}

	return nil
}

func (zone *Zone) PlaybookFile(playbookName string) string {
	if playbookName == "" {
		playbookName = DefaultPlaybook
	}

	if strings.HasSuffix(playbookName, ".yaml") {
		return path.Join(zone.Product.dir, playbookName)
	}

	f := path.Join(zone.Product.dir, playbookName+".yaml")
	if _, err := os.Stat(f); !os.IsNotExist(err) {
		return f
	}

	return path.Join(zone.Product.dir, DefaultPlaybookFile)
}

func (zone *Zone) SetComponentVersion(componentName string, componentVersion string) error {
	if !zone.Product.HasComponent(componentName) {
		msg := fmt.Sprintf("zone does not have component: (%s)", componentName)
		return errors.New(msg)
	}
	zone.Product.Components[componentName].Version = componentVersion
	return nil
}

func (zone *Zone) LoadHosts() error {
	b, err := os.ReadFile(zone.HostsFile)
	if err != nil {
		msg := fmt.Sprintf("read file (%s) failed, err: %s", zone.HostsFile, err)
		return errors.New(msg)
	}

	i := ansible.NewAnsibleInventory()
	if err := yaml.Unmarshal(b, i); err != nil {
		msg := fmt.Sprintf("unmarshal hosts failed, err: %s", err)
		return errors.New(msg)
	}

	zone.CMDB.Inventory = i
	return nil
}

func (zone *Zone) GetK8SForComponent(componentName string) *K8S {
	if platform, ok := zone.CMDB.Platforms[componentName]; ok {
		if platform.K8S != nil {
			return platform.K8S
		}
	}

	if platform, ok := zone.CMDB.Platforms["all"]; ok {
		if platform.K8S != nil {
			return platform.K8S
		}
	}

	return &K8S{}
}

func (zone *Zone) GetK8SForProduct() *K8S {
	if platform, ok := zone.CMDB.Platforms["all"]; ok {
		if platform.K8S != nil {
			return platform.K8S
		}
	}

	return &K8S{}
}

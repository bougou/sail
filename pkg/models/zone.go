package models

import (
	"errors"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"log"

	"github.com/bougou/gopkg/common"
	"github.com/bougou/sail/pkg/ansible"
	"gopkg.in/yaml.v3"
)

const (
	ProductMetavar = "_sail_product"
)

type SailOption struct {
	ProductsDir string
	PackagesDir string
	TargetsDir  string

	DefaultTarget string
	DefaultZone   string
}

type Zone struct {
	ProductName  string `yaml:"-"`
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

	targetBackupDir string

	baseExecCmd string

	sailOption *SailOption
}

func NewZone(sailOption *SailOption, targetName string, zoneName string) *Zone {
	zone := &Zone{
		TargetName: targetName,
		ZoneName:   zoneName,

		TargetDir: path.Join(sailOption.TargetsDir, targetName),
		ZoneDir:   path.Join(sailOption.TargetsDir, targetName, zoneName),

		HostsFile:    path.Join(sailOption.TargetsDir, targetName, zoneName, "hosts.yml"),
		VarsFile:     path.Join(sailOption.TargetsDir, targetName, zoneName, "vars.yml"),
		ComputedFile: path.Join(sailOption.TargetsDir, targetName, zoneName, "_computed.yml"),

		ResourcesDir: path.Join(sailOption.TargetsDir, targetName, zoneName, "resources"),

		HelmDir: path.Join(sailOption.TargetsDir, targetName, zoneName, ".helm"),

		CMDB:     NewCMDB(),
		Computed: make(map[string]interface{}),

		ansibleCfgFile: path.Join(sailOption.ProductsDir, "ansible.cfg"),

		sailOption: sailOption,
	}

	return zone
}

// Load initialize the zone with specified product.
// If exists is true, means the zone alredy exist, it will try to determine the product name from zone vars file.
// If exists is false, means the zone is newly created, the zone.ProductName should already be filled.
func (zone *Zone) Load(exists bool) error {
	if exists {
		productName, err := zone.determineProduct()
		if err != nil {
			msg := fmt.Sprintf("determine product failed, err: %s", err)
			return errors.New(msg)
		}
		zone.ProductName = productName
	}

	if zone.ProductName == "" {
		return errors.New("empty product name")
	}

	product := NewProduct(zone.ProductName, zone.sailOption.ProductsDir)
	if err := product.Init(); err != nil {
		msg := fmt.Sprintf("init product failed, err: %s", err)
		return errors.New(msg)
	}

	if exists {
		if err := zone.LoadHosts(); err != nil {
			msg := fmt.Sprintf("load hosts failed, err: %s", err)
			return errors.New(msg)
		}

		if err := product.LoadZone(zone.VarsFile); err != nil {
			msg := fmt.Sprintf("load zone vars failed, err: %s", err)
			return errors.New(msg)
		}
	}

	// Todo
	// if err := product.Check(); err != nil {
	// 	msg := fmt.Sprintf("product Check failed, err: %s", err)
	// 	return errors.New(msg)
	// }

	zone.Product = product
	zone.Product.Vars[ProductMetavar] = zone.ProductName

	if err := zone.PrepareHelm(); err != nil {
		msg := fmt.Sprintf("prepare zone helm failed, err: %s", err)
		return errors.New(msg)
	}

	return nil
}

func (zone *Zone) HelmChartDir() string {
	return path.Join(zone.HelmDir, zone.ProductName)
}

func (zone *Zone) PrepareHelm() error {
	if err := zone.prepareHelmDir(); err != nil {
		return err
	}

	for _, componentName := range zone.Product.ComponentList() {
		component := zone.Product.components[componentName]
		for _, role := range component.GetRoles() {
			zone.prepareHelmTemplates(componentName, role)
			zone.prepareHelmCRDs(componentName, role)
		}
	}
	return nil
}

func (zone *Zone) prepareHelmDir() error {
	if err := os.RemoveAll(zone.HelmDir); err != nil {
		msg := fmt.Sprintf("clear helm dir failed, err: %s", err)
		return errors.New(msg)
	}

	if err := os.MkdirAll(path.Join(zone.HelmChartDir(), "templates"), os.ModePerm); err != nil {
		msg := fmt.Sprintf("create helm chart templates dir failed, err: %s", err)
		return errors.New(msg)
	}

	if err := os.MkdirAll(path.Join(zone.HelmChartDir(), "crds"), os.ModePerm); err != nil {
		msg := fmt.Sprintf("create helm chart crds dir failed, err: %s", err)
		return errors.New(msg)
	}

	if err := os.Symlink(zone.ResourcesDir, path.Join(zone.HelmChartDir(), "resources")); err != nil {
		msg := fmt.Sprintf("create resources symlink failed, err: %s", err)
		return errors.New(msg)
	}

	if err := os.Symlink(zone.VarsFile, path.Join(zone.HelmChartDir(), "values.yml")); err != nil {
		msg := fmt.Sprintf("create values.yml symlink failed, err: %s", err)
		return errors.New(msg)
	}

	if _, err := os.Stat(zone.Product.helmChartFile); err == nil {
		if err := copyFile(zone.Product.helmChartFile, path.Join(zone.HelmChartDir(), "Chart.yml")); err != nil {
			msg := fmt.Sprintf("copy Chart.yml failed, err: %s", err)
			return errors.New(msg)
		}
	}

	return nil
}

func (zone *Zone) prepareHelmTemplates(componetName string, roleName string) error {
	roleDir := path.Join(zone.Product.rolesDir, roleName)
	roleHelmTemplatesDir := path.Join(roleDir, "helm", "templates")

	helmTemplates := []string{}
	filepath.WalkDir(roleHelmTemplatesDir, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() || strings.HasPrefix(entry.Name(), ".") {
			return nil
		}
		helmTemplates = append(helmTemplates, path)
		return nil
	})

	for _, helmTemplate := range helmTemplates {
		fileBasename := path.Base(helmTemplate)
		newFileBasename := fmt.Sprintf("%s-%s", componetName, fileBasename)
		dstFile := path.Join(zone.HelmChartDir(), "templates", newFileBasename)
		if err := copyFile(helmTemplate, dstFile); err != nil {
			msg := fmt.Sprintf("copy file failed, err: %s", err)
			return errors.New(msg)
		}
	}

	return nil
}

func (zone *Zone) prepareHelmCRDs(componetName string, roleName string) error {
	roleDir := path.Join(zone.Product.rolesDir, roleName)
	roleHelmCRDsDir := path.Join(roleDir, "helm", "crds")

	helmCRDs := []string{}
	filepath.WalkDir(roleHelmCRDsDir, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() || strings.HasPrefix(entry.Name(), ".") {
			return nil
		}
		helmCRDs = append(helmCRDs, path)
		return nil
	})

	for _, helmCRD := range helmCRDs {
		fileBasename := path.Base(helmCRD)
		newFileBasename := fmt.Sprintf("%s-%s", componetName, fileBasename)
		dstFile := path.Join(zone.HelmChartDir(), "crds", newFileBasename)
		if err := copyFile(helmCRD, dstFile); err != nil {
			msg := fmt.Sprintf("copy file failed, err: %s", err)
			return errors.New(msg)
		}
	}

	return nil
}

func copyFile(src, dst string) error {
	input, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(dst, input, 0644)
	if err != nil {
		return err
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

// determineProduct derive the product name from zone vars file
func (zone *Zone) determineProduct() (string, error) {
	b, err := os.ReadFile(zone.VarsFile)
	if err != nil {
		msg := fmt.Sprintf("read zone vars file failed, err: %s", err)
		return "", errors.New(msg)
	}

	m := map[string]interface{}{}

	if err := yaml.Unmarshal(b, &m); err != nil {
		msg := fmt.Sprintf("yaml unmarshal failed, err: %s", err)
		return "", errors.New(msg)
	}

	product, ok := m[ProductMetavar]
	if !ok {
		msg := fmt.Sprintf("not found (%s) variable in vars.yml file for target/zone: (%s/%s), you have to fix that before continue", ProductMetavar, zone.TargetName, zone.ZoneName)
		return "", errors.New(msg)
	}

	p, ok := product.(string)
	if !ok {
		msg := fmt.Sprintf("the value of (%s) variable is not a string", ProductMetavar)
		return "", errors.New(msg)
	}

	return p, nil
}

func (zone *Zone) HandleCompatibity() {
	// Todo
	// Domain Specific Language (Declarative)
	// migrate.yml
}

func (zone *Zone) check() error {
	return nil
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

func (zone *Zone) PlaybookFile(playbook string) string {
	if playbook == "" {
		playbook = DefaultPlaybook
	}

	if strings.HasSuffix(playbook, ".yml") || strings.HasSuffix(playbook, ".yaml") {
		return path.Join(zone.Product.dir, playbook)
	}

	var f string
	f = path.Join(zone.Product.dir, playbook+".yml")
	if _, err := os.Stat(f); !os.IsNotExist(err) {
		return f
	}
	f = path.Join(zone.Product.dir, playbook+".yaml")
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

// Todo, construct Helm chart from the helm dir of each component of the product.
func (zone *Zone) Helm() {
}

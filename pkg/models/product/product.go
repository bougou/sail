package product

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/bougou/sail/pkg/ansible"
	"github.com/bougou/sail/pkg/models/cmdb"
	"github.com/imdario/mergo"
	"github.com/jinzhu/copier"
	"gopkg.in/yaml.v3"
)

const DefaultPlaybook string = "sail"
const DefaultPlaybookFile string = "sail.yaml"

type Product struct {
	Name string `json:"Name,omitempty"  yaml:"Name,omitempty"`

	// zone specific vars of product
	Vars map[string]interface{} `json:"Vars,omitempty" yaml:"Vars,omitempty"`
	// zone specific components of product
	Components map[string]*Component `json:"Components,omitempty" yaml:"Components,omitempty"`

	// default vars of product, it is loaded from products/<productName>/vars.yaml
	DefaultVars map[string]interface{}
	// default components of product, it is loaded from products/<productName>/{components.yaml,components/*.yaml}
	DefaultComponents map[string]Component

	// installation order for components, it will be used for auto generating sail ansible playbook
	// it will be filled by product order.yaml file.
	order []string

	baseDir        string
	Dir            string
	componentsFile string
	componentsDir  string
	varsFile       string
	runFile        string
	migrateFile    string
	orderFile      string
	RolesDir       string

	defaultPlaybook string
}

func (p *Product) Compute(cm *cmdb.CMDB) error {
	for componentName, component := range p.Components {
		if err := component.Compute(cm); err != nil {
			msg := fmt.Sprintf("compute product component (%s) failed, err: %s", componentName, err)
			return errors.New(msg)
		}
	}

	return nil
}

func NewProduct(name string, baseDir string) *Product {
	p := &Product{
		Name:       name,
		Components: make(map[string]*Component),
		Vars:       make(map[string]interface{}),

		DefaultComponents: make(map[string]Component),
		DefaultVars:       make(map[string]interface{}),
		order:             make([]string, 0),

		defaultPlaybook: DefaultPlaybook,
		baseDir:         baseDir,
		Dir:             path.Join(baseDir, name),
		varsFile:        path.Join(baseDir, name, "vars.yaml"),
		runFile:         path.Join(baseDir, name, DefaultPlaybookFile),
		componentsFile:  path.Join(baseDir, name, "components.yaml"),
		componentsDir:   path.Join(baseDir, name, "components"),
		migrateFile:     path.Join(baseDir, name, "migrate.yaml"),
		orderFile:       path.Join(baseDir, name, "order.yaml"),
		RolesDir:        path.Join(baseDir, name, "roles"),
	}

	return p
}

func (p *Product) DefaultPlaybook() string {
	return p.defaultPlaybook
}

func (p *Product) SailPlaybookFile() string {
	return p.runFile
}

// Init will init product internal fields
func (p *Product) Init() error {
	if err := p.loadDefaultVars(); err != nil {
		msg := fmt.Sprintf("load product (%s) vars from file (%s) failed, err: %s", p.Name, p.varsFile, err)
		return errors.New(msg)
	}

	if err := p.loadDefaultComponents(); err != nil {
		msg := fmt.Sprintf("load product (%s) components failed, err: %s", p.Name, err)
		return errors.New(msg)
	}

	if err := p.loadOrder(); err != nil {
		msg := fmt.Sprintf("load product (%s) order from file (%s) failed, err: %s", p.Name, p.orderFile, err)
		return errors.New(msg)
	}

	return nil
}

func (p *Product) HasComponent(name string) bool {
	_, exists := p.Components[name]
	return exists
}

func (p *Product) SetComponentEnabled(name string, flag bool) error {
	if !p.HasComponent(name) {
		msg := fmt.Sprintf("can not enable component, the product does not have this component (%s)", name)
		return errors.New(msg)
	}

	p.Components[name].Enabled = flag
	if flag {
		p.Components[name].External = false
	}

	// Todo,
	// if this component is (enabled:true or external:true), we should automatically set all its depenent components to be (enabled:true or external:true)
	return nil
}

func (p *Product) SetComponentExternalEnabled(name string, flag bool) error {
	if !p.HasComponent(name) {
		msg := fmt.Sprintf("can not enable component, the product does not have this component (%s)", name)
		return errors.New(msg)
	}

	p.Components[name].External = flag
	if flag {
		p.Components[name].Enabled = false
	}
	return nil

	// Todo,
	// if this component is (enabled:true or external:true), we should automatically set all its depenent components to be (enabled:true or external:true)
}

func (p *Product) ComponentList() []string {
	out := []string{}
	for k := range p.Components {
		out = append(out, k)
	}

	sorted := sort.StringSlice(out)
	sort.Sort(sorted)
	return sorted
}

type FilterOption func(c *Component) bool

func FilterOptionEnabled(c *Component) bool {
	return c.Enabled
}

func FilterOptionDisabled(c *Component) bool {
	return !c.Enabled
}

func FilterOptionFormPod(c *Component) bool {
	return c.Form == ComponentFormPod
}

func FilterOptionFormServer(c *Component) bool {
	return c.Form != ComponentFormPod
}

func NewFilterOptionByComponentsMap(m map[string]string) FilterOption {
	return func(c *Component) bool {
		if _, ok := m[c.Name]; ok {
			return true
		}
		return false
	}
}

func (p *Product) ComponentListWithFitlerOptionsOr(filterOptions ...FilterOption) []string {
	out := []string{}
	for componentName, component := range p.Components {
		for _, filterOption := range filterOptions {
			if filterOption(component) {
				out = append(out, componentName)
				break
			}
		}
	}

	sorted := sort.StringSlice(out)
	sort.Sort(sorted)
	return sorted
}

func (p *Product) ComponentListWithFilterOptionsAnd(filterOptions ...FilterOption) []string {
	out := []string{}
	for componentName, component := range p.Components {
		pass := false
		for _, filterOption := range filterOptions {
			if !filterOption(component) {
				pass = true
				break
			}
		}

		if !pass {
			out = append(out, componentName)
		}
	}

	sorted := sort.StringSlice(out)
	sort.Sort(sorted)
	return sorted
}

// GenSail generate the default sail.yaml ansible playbook file
func (p *Product) GenSail() (ansible.Playbook, error) {
	out := ansible.Playbook(make([]ansible.Play, 0))

	gatherFactsPlay := ansible.NewPlay("gather facts", "all")
	gatherFactsPlay.GatherFacts = true
	gatherFactsPlay.AnyErrorsFatal = false
	gatherFactsPlay.Become = false
	role := ansible.Role{Role: "always"}
	gatherFactsPlay.WithRoles(role)
	gatherFactsPlay.WithTags("gather-facts")
	out = append(out, *gatherFactsPlay)

	for _, compName := range p.order {
		c, exists := p.DefaultComponents[compName]
		if !exists {
			return nil, fmt.Errorf("component (%s) does not declared by product (%s)", compName, p.Name)
		}

		play, err := c.GenAnsiblePlay()

		// Todo, make it configurable
		play.AnyErrorsFatal = false

		if err != nil {
			msg := fmt.Sprintf("gen ansible playbook for component (%s) failed, err: %s", c.Name, err)
			return nil, errors.New(msg)
		}
		out = append(out, *play)
	}

	return out, nil
}

func (p *Product) loadDefaultVars() error {
	b, err := os.ReadFile(p.varsFile)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			msg := fmt.Sprintf("not found default vars file (%s) for product (%s)", p.varsFile, p.Name)
			return errors.New(msg)
		}
		return err
	}

	if err := yaml.Unmarshal(b, &p.DefaultVars); err != nil {
		msg := fmt.Sprintf("unmarshal vars for product (%s) failed, err: %s", p.Name, err)
		return errors.New(msg)
	}

	m := make(map[string]interface{})
	if err := copier.Copy(&m, p.DefaultVars); err != nil {
		msg := fmt.Sprintf("copy default vars failed, err: %s", err)
		return errors.New(msg)
	}

	p.Vars = m

	return nil
}

// loadDefaultComponents will fill p.DefaultComponents with product operation code,
// and copy p.DefaultComponents to p.Components
func (p *Product) loadDefaultComponents() error {
	// componentFiles will holds
	// - the components.yaml if exist
	// - components/*.yaml
	componentFiles := []string{}

	if _, err := os.Stat(p.componentsFile); err == nil {
		componentFiles = append(componentFiles, p.componentsFile)
	}

	// store all found *.yaml files under components directory
	visitFn := func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() || strings.HasPrefix(entry.Name(), ".") {
			return nil
		}
		componentFiles = append(componentFiles, path)
		return nil
	}
	filepath.WalkDir(p.componentsDir, visitFn)

	errs := []error{}
	for _, file := range componentFiles {
		if err := p.loadComponentFile(file); err != nil {
			msg := fmt.Sprintf("load component file (%s) failed, err: %s", file, err)
			errs = append(errs, errors.New(msg))
		}
	}

	if len(errs) != 0 {
		errList := []string{""}
		for _, e := range errs {
			errList = append(errList, e.Error())
		}
		errString := strings.Join(errList, "\n")
		return errors.New(errString)
	}

	// After loading all components.yaml, copy p.DefaultComponents to p.Components
	var c map[string]*Component
	if err := copier.CopyWithOption(&c, p.DefaultComponents, copier.Option{DeepCopy: true}); err != nil {
		msg := fmt.Sprintf("copy default components failed, err: %s", err)
		return errors.New(msg)
	}
	p.Components = c

	return nil
}

func (p *Product) loadComponentFile(file string) error {
	b, err := os.ReadFile(file)
	if err != nil {
		msg := fmt.Sprintf("read file failed, err: %s", err)
		return errors.New(msg)
	}

	m := make(map[string]interface{})
	if err := yaml.Unmarshal(b, &m); err != nil {
		msg := fmt.Sprintf("yaml unmarshal file (%s) failed, err: %s", file, err)
		return errors.New(msg)
	}

	for k, v := range m {
		if _, exists := p.DefaultComponents[k]; exists {
			msg := fmt.Sprintf("found duplicate component definition for component(%s)", k)
			return errors.New(msg)
		}

		c, err := newComponentFromValue(k, v)
		if err != nil {
			return err
		}

		// check component declaration
		if c.Enabled && c.External {
			return fmt.Errorf("warn: enabled and external of component (%s) can not be both true,file (%s)", k, file)
		}
		p.DefaultComponents[k] = *c
	}

	return nil
}

func (p *Product) LoadZone(zoneVarsFile string) error {
	b, err := os.ReadFile(zoneVarsFile)
	if err != nil {
		msg := fmt.Sprintf("read file failed, err: %s", err)
		return errors.New(msg)
	}

	m := map[string]interface{}{}
	if err := yaml.Unmarshal(b, &m); err != nil {
		msg := fmt.Sprintf("unmarshal vars for failed, err: %s", err)
		return errors.New(msg)
	}

	for varKey, varValue := range m {
		// varKey is not a component
		if !p.HasComponent(varKey) {
			p.Vars[varKey] = varValue
			continue
		}

		comp, err := newComponentFromValue(varKey, varValue)
		if err != nil {
			return err
		}

		// Todo, optimize the merge behaviour
		// p.Components originally stores default components of the product,
		// now we merge the component value loaded from zone vars file into it.
		if mergo.Merge(p.Components[varKey], comp, mergo.WithOverride, mergo.WithOverwriteWithEmptyValue); err != nil {
			msg := fmt.Sprintf("merge failed for component (%s) failed, err: %s", varKey, err)
			return errors.New(msg)
		}

		c := p.Components[varKey]
		if c.Enabled && c.External {
			msg := fmt.Sprintf("234: enabled and external of component can not be both true, automatically set enabled to false for component (%s)", varKey)
			fmt.Println(msg)
			c.Enabled = false
		}

	}

	return nil
}

func (p *Product) Check(cm *cmdb.CMDB) error {
	errs := []error{}
	for _, c := range p.Components {
		if err := c.Check(); err != nil {
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

	// Todo call checkPortsConflict
	if err := p.checkPortsConflict(cm); err != nil {
		errmsg := "check port conflict failed"
		errmsgs = append(errmsgs, errmsg)
	}

	msg := fmt.Sprintf("check product (%s) faield, err: %s", p.Name, strings.Join(errmsgs, "; "))
	return errors.New(msg)
}

// checkPortsConflict
// Todo
// if multiple components are installed on same hosts, the listened ports of those components may be conflicted.
func (p *Product) checkPortsConflict(cm *cmdb.CMDB) error {
	return nil
}

// loadOrder init the order field of product.
// It must loaded after loadComponents
func (p *Product) loadOrder() error {
	order := make([]string, 0)

	b, err := os.ReadFile(p.orderFile)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			// use default order
			order = append(order, p.ComponentList()...)
			p.order = dedupSliceString(order)
			return nil
		}

		return fmt.Errorf("read file (%s) failed, err: %s", p.orderFile, err)
	}

	if err := yaml.Unmarshal(b, &order); err != nil {
		return fmt.Errorf("unmarshal order for product (%s) failed, err: %s", p.Name, err)
	}

	for _, v := range order {
		if _, exists := p.DefaultComponents[v]; !exists {
			return fmt.Errorf("component (%s) in order.yaml file does not declared by product (%s)", v, p.Name)
		}
	}

	order = append(order, p.ComponentList()...)
	p.order = dedupSliceString(order)
	return nil
}

func dedupSliceString(stringSlice []string) []string {
	keys := make(map[string]bool)
	out := []string{}

	for _, v := range stringSlice {
		if _, exists := keys[v]; !exists {
			keys[v] = true
			out = append(out, v)
		}
	}
	return out
}

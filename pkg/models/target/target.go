package target

import (
	"fmt"
	"os"
	"path"

	"github.com/bougou/sail/pkg/models"
)

type Target struct {
	Name string

	sailOption *models.SailOption
	dir        string
	vars       *TargetVars
}

type TargetVars struct {
	Zones map[string]interface{}
}

func NewTargetVars() *TargetVars {
	return &TargetVars{
		Zones: make(map[string]interface{}),
	}
}

func NewTarget(sailOption *models.SailOption, name string) *Target {
	return &Target{
		Name: name,

		sailOption: sailOption,
		vars:       NewTargetVars(),
		dir:        path.Join(sailOption.TargetsDir, name),
	}
}

func (t *Target) LoadAllZones() error {
	entries, err := os.ReadDir(t.dir)
	if err != nil {
		return err
	}

	zoneNames := make([]string, 0)
	for _, entry := range entries {
		zoneName := entry.Name()
		if !entry.IsDir() {
			continue
		}

		varsFile := path.Join(t.dir, zoneName, "vars.yaml")
		if _, err := os.Stat(varsFile); err != nil {
			// not a zone dir, ignore
			continue
		}
		zoneNames = append(zoneNames, zoneName)
	}

	for _, zoneName := range zoneNames {
		if err := t.LoadZone(zoneName); err != nil {
			return fmt.Errorf("load zone (%s) failed, err: %s", zoneName, err)
		}
	}
	return nil
}

func (t *Target) LoadZone(zoneName string) error {
	zone := NewZone(t.sailOption, t.Name, zoneName)
	if err := zone.Load(); err != nil {
		return fmt.Errorf("load zone (%s) failed, err: %s", zoneName, err)
	}

	if err := zone.Compute(); err != nil {
		return fmt.Errorf("compute zone (%s) failed, err: %s", zoneName, err)
	}

	zoneV := make(map[string]interface{})
	for k, v := range zone.Product.Vars {
		zoneV[k] = v
	}
	for k, v := range zone.Product.Components {
		zoneV[k] = v
	}
	zoneV["platforms"] = zone.CMDB.Platforms
	zoneV["inventory"] = zone.CMDB.Inventory

	t.vars.Zones[zoneName] = zoneV
	return nil
}

type Common struct {
	SSHPort string
	SSHUser string
	SSHPass string

	SudoNopass bool

	InstallDir string
	DataDir    string
}

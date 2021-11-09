package target

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/bougou/gopkg/common"
	"github.com/bougou/gopkg/copy"
	"github.com/bougou/gopkg/merge"
	"github.com/bougou/sail/pkg/models/product"
	"gopkg.in/yaml.v3"
)

// HelmDirOfProduct returns the helm chart directory for the product of the zone.
// It is used when the '_sail_helm_mode' is 'product'.
func (zone *Zone) HelmDirOfProduct() string {
	return path.Join(zone.HelmDir, zone.SailProduct)
}

// HelmDirOfComponent returns the helm chart directory for specified component.
// It is used when the '_sail_helm_mode' is 'component'.
func (zone *Zone) HelmDirOfComponent(componentName string) string {
	return path.Join(zone.HelmDir, componentName)
}

// PrepareHelm prepares helm chart(s) for zone.
func (zone *Zone) PrepareHelm() error {
	podComponentsEnabled := zone.Product.ComponentListWithFilterOptionsAnd(product.FilterOptionEnabled, product.FilterOptionFormPod)
	if len(podComponentsEnabled) == 0 {
		return nil
	}

	switch zone.SailHelmMode {
	case "component":
		return zone.PrepareHelmCharts()
	case "product":
		return zone.PrepareHelmChart()
	case "":
		return nil
	default:
		return fmt.Errorf("not supported helm mode: (%s)", zone.SailHelmMode)
	}
}

// PrepareHelmCharts prepares helm charts for each component of the product.
// Each component has its own helm chart.
// <target>/<zone>/helm/<componetName>/{Chart.yaml,templates,values.yaml,...}
func (zone *Zone) PrepareHelmCharts() error {
	if err := os.MkdirAll(zone.HelmDir, os.ModePerm); err != nil {
		return fmt.Errorf("create helm dir failed, err: %s", err)
	}

	// prepare the global values.yaml even when sail_helm_mode is set to "component".
	// the global values.yaml for zone is put under `zone.HelmDir`
	productValuesFile := path.Join(zone.Product.Dir, "values.yaml")
	zoneValuesFile := path.Join(zone.HelmDir, "values.yaml")
	if _, err := os.Stat(productValuesFile); err == nil {
		if err := mergeYamlFiles(zoneValuesFile, productValuesFile); err != nil {
			return fmt.Errorf("prepare global values.yaml failed, err: %s", err)
		}
	}

	for _, componentName := range zone.Product.ComponentListWithFilterOptionsAnd(product.FilterOptionFormPod, product.FilterOptionEnabled) {
		if err := zone.prepareComponentChart(componentName); err != nil {
			return fmt.Errorf("prepare chart for component (%s) failed, err: %s", componentName, err)
		}
	}

	return nil
}

// prepareComponentChart copy chart dir of the component into zone's helm dir,
// and changes values.yaml file.
func (zone *Zone) prepareComponentChart(componentName string) error {
	component, ok := zone.Product.Components[componentName]
	if !ok {
		return fmt.Errorf("not found component (%s) in product", componentName)
	}
	roleName := component.GetRoleName()

	roleDir := path.Join(zone.Product.RolesDir, roleName)
	roleChartDir := path.Join(roleDir, "helm", componentName)
	zoneComponentChartDir := zone.HelmDirOfComponent(componentName)

	if _, err := os.Stat(roleChartDir); err != nil {
		return fmt.Errorf("access chart dir (%s) for component (%s) failed, err: %s", roleChartDir, componentName, err)
	}

	// Copy component templates dir from role dir if exists.
	roleChartTemplatesDir := path.Join(roleChartDir, "templates")
	zoneComponentChartTemplatesDir := path.Join(zoneComponentChartDir, "templates")
	if err := os.RemoveAll(zoneComponentChartTemplatesDir); err != nil {
		return fmt.Errorf("clear templates dir failed, err: %s", err)
	}
	if _, err := os.Stat(roleChartTemplatesDir); err == nil {
		if err := copy.CopyDir(roleChartTemplatesDir+"/", zoneComponentChartTemplatesDir); err != nil {
			return fmt.Errorf("copy templates dir failed, err: %s", err)
		}
	}

	// Copy component crds dir from role dir if exists..
	roleChartCRDsDir := path.Join(roleChartDir, "crds")
	zoneComponentChartCRDsDir := path.Join(zoneComponentChartDir, "crds")
	if err := os.RemoveAll(zoneComponentChartCRDsDir); err != nil {
		return fmt.Errorf("clear crds dir failed, err: %s", err)
	}
	if _, err := os.Stat(roleChartCRDsDir); err == nil {
		if err := copy.CopyDir(roleChartCRDsDir+"/", zoneComponentChartCRDsDir); err != nil {
			return fmt.Errorf("copy crds dir failed, err: %s", err)
		}
	}

	// Copy component charts dir from role dir if exists.
	roleChartChartsDir := path.Join(roleChartDir, "charts")
	zoneComponentChartChartsDir := path.Join(zoneComponentChartDir, "charts")
	if err := os.RemoveAll(zoneComponentChartChartsDir); err != nil {
		return fmt.Errorf("clear charts dir failed, err: %s", err)
	}
	if _, err := os.Stat(roleChartChartsDir); err == nil {
		if err := copy.CopyDir(roleChartChartsDir+"/", zoneComponentChartChartsDir); err != nil {
			return fmt.Errorf("copy charts dir failed, err: %s", err)
		}
	}

	// Copy or Merge values.yaml. The values.yaml file under the role dir MUST be exists.
	roleChartValuesFile := path.Join(roleChartDir, "values.yaml")
	zoneComponentValuesFile := path.Join(zoneComponentChartDir, "values.yaml")
	if _, err := os.Stat(roleChartValuesFile); err != nil {
		return fmt.Errorf("read helm chart values.yaml file (%s) failed, err: %s", roleChartValuesFile, err)
	}
	if err := mergeYamlFiles(zoneComponentValuesFile, roleChartValuesFile); err != nil {
		return fmt.Errorf("merge values.yaml for component failed, err: %s", err)
	}

	// Copy Charts.yaml
	roleChartChartFile := path.Join(roleChartDir, "Chart.yaml")
	zoneComponentChartFile := path.Join(zoneComponentChartDir, "Chart.yaml")
	if err := copy.CopyFile(roleChartChartFile, zoneComponentChartFile); err != nil {
		return fmt.Errorf("copy Chart.yaml failed, err: %s", err)
	}

	// symlink resources dir
	chartResourcesDir := path.Join(zoneComponentChartDir, "resources")
	if err := creatSymlink(zone.ResourcesDir, chartResourcesDir); err != nil {
		return fmt.Errorf("create resources symlink failed, err: %s", err)
	}

	return nil
}

// PrepareHelmChart prepares helm chart for the product.
// There will be only one chart for the product.
// <target>/<zone>/helm/<productName>/{Chart.yaml,templates,values.yaml,...}
func (zone *Zone) PrepareHelmChart() error {
	if err := os.MkdirAll(zone.HelmDirOfProduct(), os.ModePerm); err != nil {
		return fmt.Errorf("create helm chart dir for product failed, err: %s", err)
	}

	if err := zone.prepareProductChartTemplates(); err != nil {
		return fmt.Errorf("prepare chart templates for product failed, err: %s", err)
	}

	if err := zone.prepareProductChartCRDs(); err != nil {
		return fmt.Errorf("prepare chart CRDs for product failed, err: %s", err)
	}

	productChartFile := path.Join(zone.Product.Dir, "Chart.yaml")
	zoneProductChartFile := path.Join(zone.HelmDirOfProduct(), "Chart.yaml")
	if err := copy.CopyFile(productChartFile, zoneProductChartFile); err != nil {
		return fmt.Errorf("copy Chart.yaml failed, err: %s", err)
	}

	// prepare the global values.yaml file IF EXISTS.
	// the global values.yaml for zone is put under `zone.HelmDir`.
	// the global values.yaml is OPTIONAL.
	productChartValuesFile := path.Join(zone.Product.Dir, "values.yaml")
	zoneProductChartValuesFile := path.Join(zone.HelmDir, "values.yaml")
	if _, err := os.Stat(productChartValuesFile); err == nil {
		if err := mergeYamlFiles(zoneProductChartValuesFile, productChartValuesFile); err != nil {
			return fmt.Errorf("prepare global values.yaml failed, err: %s", err)
		}
	}

	// symlink resources dir
	chartResourcesDir := path.Join(zone.HelmDirOfProduct(), "resources")
	if err := creatSymlink(zone.ResourcesDir, chartResourcesDir); err != nil {
		return fmt.Errorf("create resources symlink failed, err: %s", err)
	}

	return nil
}

// prepareProductChartTemplates prepares chart templates dir for the product in the zone.
//
//   * copy global templates if exists:
//        src: products/<productDir>/templates/filename
//        dst: <target>/<zone>/helm/<productName>/templates/filename
//   * copy component level templates if exists:
//     only the components with `enabled` set to `true` and `form` set to `pod` are considered.
//        src: products/<productDir>/roles/<roleName>/templates/filename
//        dst: <target>/<zone>/helm/<productName>/templates/<componentName>-filename
//
// The component level template dst file may overrite the global template dst file.
func (zone *Zone) prepareProductChartTemplates() error {
	productChartTemplatesDir := path.Join(zone.Product.Dir, "templates")
	zoneChartTemplatesDir := path.Join(zone.HelmDirOfProduct(), "templates")
	if err := os.RemoveAll(zoneChartTemplatesDir); err != nil {
		return fmt.Errorf("clear templates dir failed, err: %s", err)
	}

	if err := os.MkdirAll(zoneChartTemplatesDir, os.ModePerm); err != nil {
		return fmt.Errorf("create templates dir for product failed, err: %s", err)
	}

	if stat, err := os.Stat(productChartTemplatesDir); err == nil && stat.IsDir() {
		if err := copy.CopyDir(productChartTemplatesDir+"/", zoneChartTemplatesDir); err != nil {
			return fmt.Errorf("copy templates dir failed, err: %s", err)
		}
	}

	for _, componentName := range zone.Product.ComponentListWithFilterOptionsAnd(product.FilterOptionFormPod, product.FilterOptionEnabled) {
		component := zone.Product.Components[componentName]

		roleName := component.GetRoleName()
		roleDir := path.Join(zone.Product.RolesDir, roleName)
		roleHelmTemplatesDir := path.Join(roleDir, "helm", "templates")

		if err := filepath.WalkDir(roleHelmTemplatesDir, func(filepath string, entry fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if entry.IsDir() || strings.HasPrefix(entry.Name(), ".") {
				return nil
			}

			fileBasename := path.Base(filepath)
			newFileBasename := fmt.Sprintf("%s-%s", componentName, fileBasename)
			dstFile := path.Join(zoneChartTemplatesDir, newFileBasename)
			if err := copy.CopyFile(filepath, dstFile); err != nil {
				return fmt.Errorf("copy file failed, err: %s", err)
			}
			return nil
		}); err != nil {
			return fmt.Errorf("traverse role hlem templates dir failed, err: %s", err)
		}
	}

	return nil
}

// prepareProductChartCRDs prepares chart crds dir for the product in the zone.
//
//   * copy global crds if exists:
//        src: products/<productDir>/crds/filename
//        dst: <target>/<zone>/helm/<productName>/crds/filename
//   * copy component level crds if exists:
//     only the components with `enabled` set to `true` and `form` set to `pod` are considered.
//        src: products/<productDir>/roles/<roleName>/crds/filename
//        dst: <target>/<zone>/helm/<productName>/crds/<componentName>-filename
//
// The component level CRD dst file may overrite the global CRD dst file.
func (zone *Zone) prepareProductChartCRDs() error {
	productChartCRDsDir := path.Join(zone.Product.Dir, "crds")
	zoneChartCRDsDir := path.Join(zone.HelmDirOfProduct(), "crds")
	if err := os.RemoveAll(zoneChartCRDsDir); err != nil {
		return fmt.Errorf("clear crds dir failed, err: %s", err)
	}

	if err := os.MkdirAll(zoneChartCRDsDir, os.ModePerm); err != nil {
		return fmt.Errorf("create crds dir for product failed, err: %s", err)
	}

	if stat, err := os.Stat(productChartCRDsDir); err == nil && stat.IsDir() {
		if err := copy.CopyDir(productChartCRDsDir+"/", zoneChartCRDsDir); err != nil {
			return fmt.Errorf("copy crds dir failed, err: %s", err)
		}
	}

	for _, componentName := range zone.Product.ComponentListWithFitlerOptionsOr(product.FilterOptionFormPod) {
		component := zone.Product.Components[componentName]

		roleName := component.GetRoleName()
		roleDir := path.Join(zone.Product.RolesDir, roleName)
		roleHelmCRDsDir := path.Join(roleDir, "helm", "crds")

		if err := filepath.WalkDir(roleHelmCRDsDir, func(filepath string, entry fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if entry.IsDir() || strings.HasPrefix(entry.Name(), ".") {
				return nil
			}

			fileBasename := path.Base(filepath)
			newFileBasename := fmt.Sprintf("%s-%s", componentName, fileBasename)
			dstFile := path.Join(zoneChartCRDsDir, newFileBasename)
			if err := copy.CopyFile(filepath, dstFile); err != nil {
				return fmt.Errorf("copy file failed, err: %s", err)
			}
			return nil
		}); err != nil {
			return fmt.Errorf("traverse role hlem crds dir failed, err: %s", err)
		}
	}

	return nil
}

// creatSymlink creates newname as a symbolic link to oldname.
// It unlinks the newname if newname already is a symbolic link, then calls os.Symlink to create the link.
func creatSymlink(oldname string, newname string) error {
	if _, err := os.Lstat(newname); err == nil {
		if err := os.Remove(newname); err != nil {
			return err
		}
	}
	if err := os.Symlink(oldname, newname); err != nil {
		return err
	}
	return nil
}

// loadMapFromYamlFile returns a map for the content of the yaml file.
func loadMapFromYamlFile(filename string) (map[string]interface{}, error) {
	out := make(map[string]interface{})

	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("read file (%s) failed, err: %s", filename, err)
	}

	if err := yaml.Unmarshal(b, &out); err != nil {
		return nil, fmt.Errorf("yaml unmarshal failed, err: %s", err)
	}

	return out, nil
}

// mergeYamlFiles merges data in srcFiles to dstFile and save it.
// The dstFile will be auto created if not exists.
// The dstFile MUST be valid yaml files if exists.
// The passed srcFiles MUST be valid yaml files.
func mergeYamlFiles(dstFile string, srcFiles ...string) error {
	if len(srcFiles) == 0 {
		return nil
	}

	dst := make(map[string]interface{})
	if _, err := os.Stat(dstFile); err == nil {
		m, err := loadMapFromYamlFile(dstFile)
		if err != nil {
			return fmt.Errorf("load dst yaml file (%s) failed, err: %s", dstFile, err)
		}
		dst = m
	}

	dstM := merge.NewMap(dst)
	for _, srcFile := range srcFiles {
		m, err := loadMapFromYamlFile(srcFile)
		if err != nil {
			return fmt.Errorf("load src yaml file (%s) failed, err: %s", srcFile, err)
		}
		if err := dstM.Merge(m); err != nil {
			return fmt.Errorf("merge src yaml file (%s) failed, err: %s", srcFile, err)
		}
	}

	b, err := common.Encode("yaml", dstM.Value())
	if err != nil {
		return fmt.Errorf("encode failed, err: %s", err)
	}

	if err := os.WriteFile(dstFile, b, 0644); err != nil {
		return fmt.Errorf("update dst file failed, err: %s", err)
	}

	return nil
}

package models

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/bougou/gopkg/copy"
)

func (zone *Zone) HelmChartDirOfProduct() string {
	return path.Join(zone.HelmDir, zone.SailProduct)
}

func (zone *Zone) HelmChartDirOfComponent(componentName string) string {
	return path.Join(zone.HelmDir, componentName)
}

func (zone *Zone) PrepareHelm() error {
	switch zone.SailHelmMode {
	case "component":
		return zone.PrepareHelmCharts()
	case "product":
		return zone.PrepareHelmChart()
	default:
		return fmt.Errorf("not supported helm mode: (%s)", zone.SailHelmMode)
	}
}

// PrepareHelmCharts prepares helm chart for each components of the product.
// Each component has its own helm chart.
func (zone *Zone) PrepareHelmCharts() error {
	// <target>/<zone>/.helm/<componetName>/{Chart.yaml,templates,values.yaml}
	if err := os.RemoveAll(zone.HelmDir); err != nil {
		msg := fmt.Sprintf("clear helm dir failed, err: %s", err)
		return errors.New(msg)
	}

	if err := os.MkdirAll(zone.HelmDir, os.ModePerm); err != nil {
		msg := fmt.Sprintf("create helm dir failed, err: %s", err)
		return errors.New(msg)
	}

	for _, componentName := range zone.Product.ComponentList() {
		component := zone.Product.Components[componentName]
		if component.Form == ComponentFormPod {
			roleName := component.GetRoleName()

			if err := zone.prepareComponentChart(componentName, roleName); err != nil {
				return fmt.Errorf("prepare chart for component (%s) / role (%s) failed, err: %s", componentName, roleName, err)
			}
		}
	}
	return nil
}

func (zone *Zone) prepareComponentChart(componentName string, roleName string) error {
	roleDir := path.Join(zone.Product.rolesDir, roleName)
	helmChartDir := path.Join(roleDir, "helm", componentName)

	_, err := os.Stat(helmChartDir)
	if err != nil && os.IsNotExist(err) {
		return nil
	}

	if err := copy.CopyDir(helmChartDir, zone.HelmDir); err != nil {
		return fmt.Errorf("copy chart dir failed, err: %s", err)
	}

	return nil
}

// PrepareHelmChart prepares helm chart for the product.
// There will be only one chart for the product.
func (zone *Zone) PrepareHelmChart() error {
	if err := os.RemoveAll(zone.HelmDir); err != nil {
		msg := fmt.Sprintf("clear helm dir failed, err: %s", err)
		return errors.New(msg)
	}

	if err := os.MkdirAll(path.Join(zone.HelmChartDirOfProduct(), "templates"), os.ModePerm); err != nil {
		msg := fmt.Sprintf("create helm chart templates dir failed, err: %s", err)
		return errors.New(msg)
	}

	if err := os.MkdirAll(path.Join(zone.HelmChartDirOfProduct(), "crds"), os.ModePerm); err != nil {
		msg := fmt.Sprintf("create helm chart crds dir failed, err: %s", err)
		return errors.New(msg)
	}

	if err := os.Symlink(zone.ResourcesDir, path.Join(zone.HelmChartDirOfProduct(), "resources")); err != nil {
		msg := fmt.Sprintf("create resources symlink failed, err: %s", err)
		return errors.New(msg)
	}

	if err := os.Symlink(zone.VarsFile, path.Join(zone.HelmChartDirOfProduct(), "values.yaml")); err != nil {
		msg := fmt.Sprintf("create values.yaml symlink failed, err: %s", err)
		return errors.New(msg)
	}

	if _, err := os.Stat(zone.Product.helmChartFile); err == nil {
		if err := copy.CopyFile(zone.Product.helmChartFile, path.Join(zone.HelmChartDirOfProduct(), "Chart.yaml")); err != nil {
			msg := fmt.Sprintf("copy Chart.yaml failed, err: %s", err)
			return errors.New(msg)
		}
	}

	for _, componentName := range zone.Product.ComponentList() {
		component := zone.Product.Components[componentName]
		if component.Form == ComponentFormPod {
			roleName := component.GetRoleName()
			zone.prepareHelmTemplates(componentName, roleName)
			zone.prepareHelmCRDs(componentName, roleName)
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
		dstFile := path.Join(zone.HelmChartDirOfProduct(), "templates", newFileBasename)
		if err := copy.CopyFile(helmTemplate, dstFile); err != nil {
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
		dstFile := path.Join(zone.HelmChartDirOfProduct(), "crds", newFileBasename)
		if err := copy.CopyFile(helmCRD, dstFile); err != nil {
			msg := fmt.Sprintf("copy file failed, err: %s", err)
			return errors.New(msg)
		}
	}

	return nil
}

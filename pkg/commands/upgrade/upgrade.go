package upgrade

import (
	"errors"
	"fmt"
	"strings"

	"github.com/bougou/gopkg/common"
	"github.com/bougou/sail/pkg/models"
	"github.com/bougou/sail/pkg/options"
	"github.com/spf13/cobra"
)

func NewCmdUpgrade(sailOption *models.SailOption) *cobra.Command {
	o := NewUpgradeOptions(sailOption)

	cmd := &cobra.Command{
		Use:   "upgrade",
		Short: "upgrade",
		Long:  "upgrade",
		Run: func(cmd *cobra.Command, args []string) {
			common.CheckErr(o.Run(args))
		},
	}

	cmd.Flags().StringVarP(&o.TargetName, "target", "t", o.TargetName, "target name")
	cmd.MarkFlagRequired("target")

	cmd.Flags().StringVarP(&o.ZoneName, "zone", "z", o.ZoneName, "zone name")
	cmd.MarkFlagRequired("zone")

	cmd.Flags().StringArrayVarP(&o.Components, "component", "c", o.Components, "the component")
	return cmd
}

type UpgradeOptions struct {
	TargetName string `json:"target_name"`
	ZoneName   string `json:"zone_name"`

	Components []string `json:"component"`

	sailOption *models.SailOption
}

func NewUpgradeOptions(sailOption *models.SailOption) *UpgradeOptions {
	return &UpgradeOptions{
		Components: make([]string, 0),
		sailOption: sailOption,
	}
}

func (o *UpgradeOptions) Complete() error {

	return nil
}

func (o *UpgradeOptions) Validate() error {

	return nil
}

func (o *UpgradeOptions) Run(args []string) error {
	zone := models.NewZone(o.sailOption, o.TargetName, o.ZoneName)
	if err := zone.Load(true); err != nil {
		return err
	}

	m, err := options.ParseComponentOption(o.Components)
	if err != nil {
		msg := fmt.Sprintf("ParseComponentOption failed, err: %s", err)
		return errors.New(msg)
	}

	ansiblePlaybookTags := []string{}
	for componentName, componentVersion := range m {
		tag := "update-" + componentName
		ansiblePlaybookTags = append(ansiblePlaybookTags, tag)

		if componentVersion == "" {
			continue
		}
		zone.SetComponentVersion(componentName, componentVersion)
	}

	if err := zone.Dump(); err != nil {
		msg := fmt.Sprintf("zone.Dump failed, err: %s", err)
		return errors.New(msg)
	}

	rz := models.NewRunningZone(zone, zone.Product.DefaultPlaybook())

	playbookArgs := []string{}
	if len(ansiblePlaybookTags) != 0 {
		// Note, CANNOT pass "--tags tag1,tag2" as one item into the slice
		playbookArgs = append(playbookArgs, "--tags", strings.Join(ansiblePlaybookTags, ","))
	}

	playbookArgs = append(playbookArgs, args...)

	return rz.Run(playbookArgs)
}

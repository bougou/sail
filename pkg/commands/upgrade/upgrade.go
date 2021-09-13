package upgrade

import (
	"errors"
	"fmt"

	"github.com/bougou/gopkg/common"
	"github.com/bougou/sail/pkg/models"
	"github.com/bougou/sail/pkg/options"
	"github.com/spf13/cobra"
)

func NewCmdUpgrade(sailOption *models.SailOption) *cobra.Command {
	o := NewUpgradeOptions(sailOption)

	cmd := &cobra.Command{
		Use:   "upgrade",
		Short: "upgrade the components of product to specified version",
		Long:  "upgrade the components of product to specified version",
		Run: func(cmd *cobra.Command, args []string) {
			common.CheckErr(o.Complete(cmd, args))
			common.CheckErr(o.Validate())
			common.CheckErr(o.Run(args))
		},
	}

	cmd.Flags().StringVarP(&o.TargetName, "target", "t", o.TargetName, "target name")
	cmd.Flags().StringVarP(&o.ZoneName, "zone", "z", o.ZoneName, "zone name")
	cmd.Flags().StringArrayVarP(&o.Components, "component", "c", o.Components, "the component")
	cmd.Flags().BoolVarP(&o.Ansible, "ansible", "", o.Ansible, "choose components deployed as server")
	cmd.Flags().BoolVarP(&o.Helm, "helm", "", o.Helm, "choose components deployed as pod")
	return cmd
}

type UpgradeOptions struct {
	TargetName string `json:"target_name"`
	ZoneName   string `json:"zone_name"`

	Components []string `json:"component"`
	Ansible    bool     `json:"ansible"`
	Helm       bool     `json:"helm"`

	sailOption *models.SailOption
}

func NewUpgradeOptions(sailOption *models.SailOption) *UpgradeOptions {
	return &UpgradeOptions{
		Components: make([]string, 0),
		sailOption: sailOption,
	}
}

func (o *UpgradeOptions) Complete(cmd *cobra.Command, args []string) error {
	if o.TargetName == "" {
		o.TargetName = o.sailOption.DefaultTarget
	}
	if o.ZoneName == "" {
		o.ZoneName = o.sailOption.DefaultZone
	}

	return nil
}

func (o *UpgradeOptions) Validate() error {
	if o.TargetName == "" {
		return errors.New("must specify target name")
	}
	if o.ZoneName == "" {
		return errors.New("must specify zone name")
	}
	return nil
}

func (o *UpgradeOptions) Run(args []string) error {
	fmt.Printf("👉 target: (%s), zone: (%s)\n", o.TargetName, o.ZoneName)
	zone := models.NewZone(o.sailOption, o.TargetName, o.ZoneName)
	if err := zone.Load(); err != nil {
		return err
	}

	serverComponents, podComponents, err := options.ParseChoosedComponents(zone, o.Components, o.Ansible, o.Helm)
	if err != nil {
		msg := fmt.Sprintf("parse component option failed, err: %s", err)
		return errors.New(msg)
	}

	if err := zone.Dump(); err != nil {
		msg := fmt.Sprintf("zone.Dump failed, err: %s", err)
		return errors.New(msg)
	}

	var ansiblePlaybookTags []string
	for componentName := range serverComponents {
		// Note: Ansible Tag for update component
		ansiblePlaybookTag := "update-" + componentName
		ansiblePlaybookTags = append(ansiblePlaybookTags, ansiblePlaybookTag)
	}

	rz := models.NewRunningZone(zone, zone.Product.DefaultPlaybook())
	rz.WithServerComponents(serverComponents)
	rz.WithPodComponents(podComponents)
	rz.WithAnsiblePlaybookTags(ansiblePlaybookTags)

	return rz.Run(args)
}

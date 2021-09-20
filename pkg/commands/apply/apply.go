package apply

import (
	"errors"
	"fmt"

	"github.com/bougou/gopkg/common"
	"github.com/bougou/sail/pkg/models"
	"github.com/bougou/sail/pkg/models/target"
	"github.com/bougou/sail/pkg/options"
	"github.com/spf13/cobra"
)

func NewCmdApply(sailOption *models.SailOption) *cobra.Command {
	o := NewApplyOptions(sailOption)

	cmd := &cobra.Command{
		Use:   "apply",
		Short: "apply start deployment",
		Long:  "apply start deployment",
		Run: func(cmd *cobra.Command, args []string) {
			// cmd here refers to parent command
			common.CheckErr(o.Complete(cmd, args))
			common.CheckErr(o.Validate())
			common.CheckErr(o.Run(args))
		},
	}

	cmd.Flags().StringVarP(&o.TargetName, "target", "t", o.TargetName, "target name")
	cmd.Flags().StringVarP(&o.ZoneName, "zone", "z", o.ZoneName, "zone name")
	cmd.Flags().StringVarP(&o.Playbook, "playbook", "p", "run", "optional playbook name")
	cmd.Flags().StringVarP(&o.StartAtPlay, "start-at-play", "", "", "start the playbook from the play with this tag name")
	cmd.Flags().StringArrayVarP(&o.Components, "component", "c", o.Components, "the component")
	cmd.Flags().BoolVarP(&o.Ansible, "ansible", "", o.Ansible, "choose components deployed as server")
	cmd.Flags().BoolVarP(&o.Helm, "helm", "", o.Helm, "choose components deployed as pod")

	return cmd
}

type ApplyOptions struct {
	TargetName string `json:"target_name"`
	ZoneName   string `json:"zone_name"`
	Playbook   string `json:"playbook"`

	StartAtPlay string `json:"start_at_playbook"`

	Components []string `json:"component"`
	Ansible    bool     `json:"ansible"`
	Helm       bool     `json:"helm"`

	sailOption *models.SailOption
}

func NewApplyOptions(sailOption *models.SailOption) *ApplyOptions {
	return &ApplyOptions{
		sailOption: sailOption,
	}
}

func (o *ApplyOptions) Complete(cmd *cobra.Command, args []string) error {
	if o.TargetName == "" {
		o.TargetName = o.sailOption.DefaultTarget
	}
	if o.ZoneName == "" {
		o.ZoneName = o.sailOption.DefaultZone
	}

	return nil
}

func (o *ApplyOptions) Validate() error {
	return nil
}

func (o *ApplyOptions) Run(args []string) error {
	options.PrintColorHeader(o.TargetName, o.ZoneName)

	zone := target.NewZone(o.sailOption, o.TargetName, o.ZoneName)
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
		ansiblePlaybookTag := componentName
		ansiblePlaybookTags = append(ansiblePlaybookTags, ansiblePlaybookTag)
	}

	rz := target.NewRunningZone(zone, o.Playbook)
	rz.WithServerComponents(serverComponents)
	rz.WithPodComponents(podComponents)
	rz.WithAnsiblePlaybookTags(ansiblePlaybookTags)
	rz.WithStartAtPlay(o.StartAtPlay)

	return rz.Run(args)
}

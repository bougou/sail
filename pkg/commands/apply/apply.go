package apply

import (
	"errors"
	"fmt"
	"strings"

	"github.com/bougou/gopkg/common"
	"github.com/bougou/sail/pkg/ansible"
	"github.com/bougou/sail/pkg/models"
	"github.com/spf13/cobra"
)

func NewCmdApply(sailOption *models.SailOption) *cobra.Command {
	o := NewApplyOptions(sailOption)

	cmd := &cobra.Command{
		Use:   "apply",
		Short: "apply",
		Long:  "apply",
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

	return cmd
}

type ApplyOptions struct {
	TargetName string `json:"target_name"`
	ZoneName   string `json:"zone_name"`
	Playbook   string `json:"playbook"`

	StartAtPlay string `json:"start_at_playbook"`
	sailOption  *models.SailOption
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
	fmt.Printf("👉 target: (%s), zone: (%s)\n", o.TargetName, o.ZoneName)
	zone := models.NewZone(o.sailOption, o.TargetName, o.ZoneName)
	if err := zone.Load(true); err != nil {
		return err
	}

	if err := zone.Dump(); err != nil {
		msg := fmt.Sprintf("zone.Dump failed, err: %s", err)
		return errors.New(msg)
	}

	playbookFile := zone.PlaybookFile(o.Playbook)
	playbook, err := ansible.NewPlaybookFromFile(playbookFile)
	if err != nil {
		return err
	}

	playbookTags := []string{}
	if o.StartAtPlay != "" {
		playbookTags = playbook.PlaysTagsStartAt(o.StartAtPlay)
	}

	playbookArgs := []string{}
	if len(playbookTags) != 0 {
		// Note, CANNOT pass "--tags tag1,tag2" as one item into the slice
		playbookArgs = append(playbookArgs, "--tags", strings.Join(playbookTags, ","))
	}

	playbookArgs = append(playbookArgs, args...)

	rz := models.NewRunningZone(zone, o.Playbook)
	return rz.Run(playbookArgs)
}

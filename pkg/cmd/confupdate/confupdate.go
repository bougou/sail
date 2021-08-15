package confupdate

import (
	"errors"
	"fmt"

	"github.com/bougou/sail/pkg/options"
	cmdutil "github.com/bougou/sail/pkg/util"

	"github.com/bougou/sail/pkg/models"
	"github.com/spf13/cobra"
)

func NewCmdConfUpdate(sailOption *models.SailOption) *cobra.Command {
	o := NewConfUpdateOptions(sailOption)

	cmd := &cobra.Command{
		Use:   "conf-update",
		Short: "conf-update",
		Long:  "conf-update",
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(o.Run())
		},
	}

	cmd.Flags().StringVarP(&o.TargetName, "target", "t", o.TargetName, "target name")
	cmd.MarkFlagRequired("target")

	cmd.Flags().StringVarP(&o.ZoneName, "zone", "z", o.ZoneName, "zone name")
	cmd.MarkFlagRequired("zone")

	cmd.Flags().StringArrayVarP(&o.Hosts, "hosts", "", nil, "the hosts")
	cmd.Flags().StringArrayVarP(&o.Components, "components", "c", nil, "enable components")
	cmd.Flags().StringArrayVarP(&o.NoComponents, "no-components", "", nil, "disable components")
	cmd.Flags().StringArrayVarP(&o.ExternalComponents, "external-components", "", nil, "enable external components")
	cmd.Flags().StringArrayVarP(&o.NoExternalComponents, "no-external-components", "", nil, "disable external components")

	return cmd
}

type ConfUpdateOptions struct {
	TargetName string
	ZoneName   string

	Hosts []string

	Components           []string
	NoComponents         []string
	ExternalComponents   []string
	NoExternalComponents []string

	sailOption *models.SailOption
}

func NewConfUpdateOptions(sailOption *models.SailOption) *ConfUpdateOptions {
	return &ConfUpdateOptions{
		sailOption: sailOption,
	}
}

func (o *ConfUpdateOptions) Complete() error {

	return nil
}

func (o *ConfUpdateOptions) Validate() error {

	return nil
}

func (o *ConfUpdateOptions) Run() error {
	zone := models.NewZone(o.sailOption, o.TargetName, o.ZoneName)
	if err := zone.Load(true); err != nil {
		msg := fmt.Sprintf("zone.Load failed, err: %s", err)
		return errors.New(msg)
	}

	m, err := options.ParseHostsOption(o.Hosts)
	if err != nil {
		msg := fmt.Sprintf("ParseHostsOption failed, err: %s", err)
		return errors.New(msg)
	}
	zone.PatchActionHostsMap(m)

	for _, c := range options.ParseComponentsOption(o.Components) {
		zone.Product.SetComponentEnabled(c, true)
	}

	for _, c := range options.ParseComponentsOption(o.NoComponents) {
		zone.Product.SetComponentEnabled(c, false)
	}

	for _, c := range options.ParseComponentsOption(o.ExternalComponents) {
		zone.Product.SetComponentExternalEnabled(c, true)
	}

	for _, c := range options.ParseComponentsOption(o.NoExternalComponents) {
		zone.Product.SetComponentExternalEnabled(c, false)
	}

	if err := zone.Dump(); err != nil {
		msg := fmt.Sprintf("zone.Dump failed, err: %s", err)
		return errors.New(msg)
	}

	return nil
}

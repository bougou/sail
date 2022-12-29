package confupdate

import (
	"errors"
	"fmt"

	"github.com/bougou/gopkg/common"
	"github.com/bougou/sail/pkg/models"
	"github.com/bougou/sail/pkg/models/target"
	"github.com/bougou/sail/pkg/options"
	"github.com/spf13/cobra"
)

func NewCmdConfUpdate(sailOption *models.SailOption) *cobra.Command {
	o := NewConfUpdateOptions(sailOption)

	cmd := &cobra.Command{
		Use:   "conf-update",
		Short: "update the vars for an environment",
		Long:  "update the vars for an environment",
		Run: func(cmd *cobra.Command, args []string) {
			common.CheckErr(o.Complete(cmd, args))
			common.CheckErr(o.Validate())
			common.CheckErr(o.Run())
		},
	}

	cmd.Flags().StringVarP(&o.TargetName, "target", "t", o.TargetName, "target name")

	cmd.Flags().StringVarP(&o.ZoneName, "zone", "z", o.ZoneName, "zone name")

	cmd.Flags().StringArrayVarP(&o.Hosts, "hosts", "", nil, "host string, eg: --hosts A,B/10.0.0.1,10.0.0.2 --hosts +C/10.0.0.3,10.0.0.4 --hosts -C,D,E/10.0.0.4 --hosts 10.0.0.1")
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

func (o *ConfUpdateOptions) Complete(cmd *cobra.Command, args []string) error {
	if o.TargetName == "" {
		o.TargetName = o.sailOption.DefaultTarget
	}
	if o.ZoneName == "" {
		o.ZoneName = o.sailOption.DefaultZone
	}

	return nil
}

func (o *ConfUpdateOptions) Validate() error {
	if o.TargetName == "" {
		return errors.New("must specify target name")
	}
	if o.ZoneName == "" {
		return errors.New("must specify zone name")
	}

	return nil
}

func (o *ConfUpdateOptions) Run() error {
	options.PrintColorHeader(o.TargetName, o.ZoneName)

	zone := target.NewZone(o.sailOption, o.TargetName, o.ZoneName)
	if err := zone.Load(); err != nil {
		return fmt.Errorf("zone.Load failed, err: %s", err)
	}

	m, err := options.ParseHostsOptions(o.Hosts)
	if err != nil {
		return fmt.Errorf("parse hosts option failed, err: %s", err)
	}
	if err := zone.PatchActionHostsMap(m); err != nil {
		return fmt.Errorf("patch hosts failed, err: %s", err)
	}

	if components, err := options.ParseComponentsOption(o.Components); err != nil {
		return fmt.Errorf("parse component options failed, err: %s", err)
	} else {
		for c := range components {
			if err := zone.Product.SetComponentEnabled(c, true); err != nil {
				return fmt.Errorf("update component enabled to true failed, err: %s", err)
			}
		}
	}

	if components, err := options.ParseComponentsOption(o.NoComponents); err != nil {
		return fmt.Errorf("parse component options failed, err: %s", err)
	} else {
		for c := range components {
			if err := zone.Product.SetComponentEnabled(c, false); err != nil {
				return fmt.Errorf("update component enabled to false failed, err: %s", err)
			}
		}
	}

	if components, err := options.ParseComponentsOption(o.ExternalComponents); err != nil {
		return fmt.Errorf("parse component options failed, err: %s", err)
	} else {
		for c := range components {
			if err := zone.Product.SetComponentExternalEnabled(c, true); err != nil {
				return fmt.Errorf("update component external to true failed, err: %s", err)
			}
		}
	}

	if components, err := options.ParseComponentsOption(o.NoExternalComponents); err != nil {
		return fmt.Errorf("parse component options failed, err: %s", err)
	} else {
		for c := range components {
			if err := zone.Product.SetComponentExternalEnabled(c, false); err != nil {
				return fmt.Errorf("update component external to false failed, err: %s", err)
			}
		}
	}

	if err := zone.Dump(); err != nil {
		return fmt.Errorf("zone.Dump failed, err: %s", err)
	}

	return nil
}

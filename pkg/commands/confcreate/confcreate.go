package confcreate

import (
	"errors"
	"fmt"
	"os"

	"github.com/bougou/gopkg/common"
	"github.com/bougou/sail/pkg/models"
	"github.com/bougou/sail/pkg/options"
	"github.com/spf13/cobra"
)

var (
	defaultInstallDir = "/opt"
	defaultDataDir    = "/data"
	defaultSSHUser    = "root"
	defaultSSHPort    = 22
)

func NewCmdConfCreate(sailOption *models.SailOption) *cobra.Command {
	o := NewConfCreateOptions(sailOption)

	cmd := &cobra.Command{
		Use:   "conf-create",
		Short: "create a new envinronment (target/zone)",
		Long:  "create a new envinronment (target/zone)",
		Run: func(cmd *cobra.Command, args []string) {
			common.CheckErr(o.Complete(cmd, args))
			common.CheckErr(o.Validate())
			common.CheckErr(o.Run())
		},
	}

	cmd.Flags().StringVarP(&o.TargetName, "target", "t", o.TargetName, "the target name")
	cmd.MarkFlagRequired("target")
	cmd.Flags().StringVarP(&o.ZoneName, "zone", "z", o.ZoneName, "the zone name")
	cmd.MarkFlagRequired("zone")
	cmd.Flags().StringVarP(&o.ProductName, "product", "p", o.ProductName, "the product name")
	cmd.MarkFlagRequired("product")

	cmd.Flags().StringVar(&o.InstallDir, "install-dir", defaultInstallDir, "the install dir")
	cmd.Flags().StringVar(&o.DataDir, "data-dir", defaultDataDir, "the data dir")
	cmd.Flags().StringVar(&o.SSHUser, "ssh-user", defaultSSHUser, "the ssh user")
	cmd.Flags().IntVar(&o.SSHPort, "ssh-port", defaultSSHPort, "the ssh port")

	cmd.Flags().StringArrayVarP(&o.Hosts, "hosts", "", o.Hosts, "the hosts")

	return cmd
}

type ConfCreateOptions struct {
	TargetName  string `json:"target_name"`
	ZoneName    string `json:"zone_name"`
	ProductName string `json:"product_name"`

	InstallDir string `json:"install_dir"`
	DataDir    string `json:"data_dir"`
	SSHUser    string `json:"ssh_user"`
	SSHPort    int    `json:"ssh_port"`

	Hosts []string `json:"hosts"`

	sailOption *models.SailOption
}

func NewConfCreateOptions(sailOption *models.SailOption) *ConfCreateOptions {
	return &ConfCreateOptions{
		sailOption: sailOption,
	}
}

func (o *ConfCreateOptions) Complete(cmd *cobra.Command, args []string) error {
	return nil
}

func (o *ConfCreateOptions) Validate() error {

	if len(o.Hosts) == 0 {
		return fmt.Errorf("must specify at least one --hosts option when create a target/zone")
	}
	return nil
}

func (o *ConfCreateOptions) Run() error {
	fmt.Printf("👉 target: (%s), zone: (%s)\n", o.TargetName, o.ZoneName)

	zone := models.NewZone(o.sailOption, o.TargetName, o.ZoneName)
	if _, err := os.Stat(zone.ZoneDir); !os.IsNotExist(err) {
		msg := fmt.Sprintf("target/zone (%s/%s) already exists, found zone dir: %s, remove the dir if you want to recreate the zone", o.TargetName, o.ZoneName, zone.ZoneDir)
		return errors.New(msg)
	}

	zone.ZoneMeta = &models.ZoneMeta{
		SailProduct:  o.ProductName,
		SailHelmMode: "component",
	}

	if err := zone.LoadNew(); err != nil {
		msg := fmt.Sprintf("zone.Load failed, err: %s", err)
		return errors.New(msg)
	}

	m, err := options.ParseHostsOption(o.Hosts)
	if err != nil {
		msg := fmt.Sprintf("parse hosts option failed, err: %s", err)
		return errors.New(msg)
	}

	if err := zone.PatchActionHostsMap(m); err != nil {
		return err
	}

	if err := zone.Dump(); err != nil {
		msg := fmt.Sprintf("dump zone failed, err: %s", err)
		return errors.New(msg)
	}

	return nil

}

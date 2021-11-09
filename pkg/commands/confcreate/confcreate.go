package confcreate

import (
	"fmt"
	"os"

	"github.com/bougou/gopkg/common"
	"github.com/bougou/sail/pkg/models"
	"github.com/bougou/sail/pkg/models/cmdb"
	"github.com/bougou/sail/pkg/models/target"
	"github.com/bougou/sail/pkg/options"
	"github.com/spf13/cobra"
)

var (
	defaultInstallDir  = "/opt"
	defaultDataDir     = "/data"
	defaultSSHUser     = "root"
	defaultSSHPort     = 22
	defaultKubeConfig  = "~/.kube/config"
	defaultKubeContext = ""
	defaultNamespace   = "default"
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
	_ = cmd.MarkFlagRequired("target")
	cmd.Flags().StringVarP(&o.ZoneName, "zone", "z", o.ZoneName, "the zone name")
	_ = cmd.MarkFlagRequired("zone")
	cmd.Flags().StringVarP(&o.ProductName, "product", "p", o.ProductName, "the product name")
	_ = cmd.MarkFlagRequired("product")

	cmd.Flags().StringVar(&o.InstallDir, "install-dir", defaultInstallDir, "the install dir")
	cmd.Flags().StringVar(&o.DataDir, "data-dir", defaultDataDir, "the data dir")
	cmd.Flags().StringVar(&o.SSHUser, "ssh-user", defaultSSHUser, "the ssh user")
	cmd.Flags().IntVar(&o.SSHPort, "ssh-port", defaultSSHPort, "the ssh port")

	cmd.Flags().StringArrayVarP(&o.Hosts, "hosts", "", o.Hosts, "the hosts")

	cmd.Flags().StringVar(&o.KubeConfig, "kubeconfig", defaultKubeConfig, "path to the kubeconfig file")
	cmd.Flags().StringVar(&o.KubeContext, "kube-context", defaultKubeContext, "name of the kubeconfig context to use")
	cmd.Flags().StringVar(&o.Namespace, "namespace", defaultNamespace, "k8s namespace scope")

	return cmd
}

type ConfCreateOptions struct {
	TargetName  string
	ZoneName    string
	ProductName string

	InstallDir string
	DataDir    string
	SSHUser    string
	SSHPort    int

	Hosts []string

	KubeConfig  string
	KubeContext string
	Namespace   string

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
	options.PrintColorHeader(o.TargetName, o.ZoneName)

	zone := target.NewZone(o.sailOption, o.TargetName, o.ZoneName)
	if _, err := os.Stat(zone.ZoneDir); !os.IsNotExist(err) {
		return fmt.Errorf("target/zone (%s/%s) already exists, found zone dir: %s, remove the dir if you want to recreate the zone", o.TargetName, o.ZoneName, zone.ZoneDir)
	}

	zone.ZoneMeta = &target.ZoneMeta{
		SailProduct:  o.ProductName,
		SailHelmMode: "component",
	}

	if err := zone.LoadNew(); err != nil {
		return fmt.Errorf("zone.Load failed, err: %s", err)
	}

	m, err := options.ParseHostsOptions(o.Hosts)
	if err != nil {
		return fmt.Errorf("parse hosts option failed, err: %s", err)
	}

	platform := cmdb.Platform{
		K8S: &cmdb.K8S{
			KubeConfig:  o.KubeConfig,
			KubeContext: o.KubeContext,
			Namespace:   o.Namespace,
		},
	}
	zone.CMDB.Platforms["all"] = platform

	if err := zone.PatchActionHostsMap(m); err != nil {
		return err
	}

	if err := zone.Dump(); err != nil {
		return fmt.Errorf("dump zone failed, err: %s", err)
	}

	return nil

}

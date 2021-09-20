package gencert

import (
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/bougou/gopkg/cert"
	"github.com/bougou/gopkg/common"
	"github.com/bougou/sail/pkg/models"
	"github.com/bougou/sail/pkg/models/target"
	"github.com/bougou/sail/pkg/options"
	"github.com/spf13/cobra"
)

const CertValidDays int = 999 * 365

func NewCmdGenCert(sailOption *models.SailOption) *cobra.Command {
	o := NewGenCertOptions(sailOption)

	cmd := &cobra.Command{
		Use:   "gen-cert",
		Short: "generate private key and certificate",
		Long:  "generate private key and certificate",
		Run: func(cmd *cobra.Command, args []string) {
			common.CheckErr(o.Complete(cmd, args))
			common.CheckErr(o.Validate())
			common.CheckErr(o.Run(args))
		},
	}

	cmd.Flags().StringVarP(&o.TargetName, "target", "t", o.TargetName, "target name")
	cmd.Flags().StringVarP(&o.ZoneName, "zone", "z", o.ZoneName, "zone name")

	cmd.Flags().StringVarP(&o.OutputDir, "output-dir", "o", o.OutputDir, "the output dir")

	cmd.Flags().StringVarP(&o.CAName, "ca", "", o.CAName, "CA name")
	cmd.MarkFlagRequired("ca")
	cmd.Flags().StringArrayVarP(&o.Names, "name", "n", o.Names, "the commonName for the cert")
	return cmd
}

type GenCertOptions struct {
	TargetName string `json:"target_name"`
	ZoneName   string `json:"zone_name"`

	OutputDir string   `json:"output_dir"`
	CAName    string   `json:"ca_name"`
	Names     []string `json:"names"`

	sailOption *models.SailOption
}

func NewGenCertOptions(sailOption *models.SailOption) *GenCertOptions {
	return &GenCertOptions{
		Names:      make([]string, 0),
		sailOption: sailOption,
	}
}

func (o *GenCertOptions) Complete(cmd *cobra.Command, args []string) error {
	if o.TargetName == "" {
		o.TargetName = o.sailOption.DefaultTarget
	}
	if o.ZoneName == "" {
		o.ZoneName = o.sailOption.DefaultZone
	}

	if o.OutputDir == "" {
		if o.TargetName == "" || o.ZoneName == "" {
			msg := fmt.Sprintf("must specify --output-dir or specify target/zone name")
			return errors.New(msg)
		}

		zone := target.NewZone(o.sailOption, o.TargetName, o.ZoneName)
		o.OutputDir = path.Join(zone.ZoneDir, "resources", "certs")
	}
	return nil
}

func (o *GenCertOptions) Validate() error {
	return nil
}

func (o *GenCertOptions) Run(args []string) error {
	fmt.Printf("output dir for certs: %s\n", o.OutputDir)
	if o.TargetName != "" && o.ZoneName != "" {
		options.PrintColorHeader(o.TargetName, o.ZoneName)
	}

	if o.CAName != "" && len(o.Names) == 0 {
		return o.genCA()
	}

	return o.genCert()
}

func (o *GenCertOptions) genCA() error {
	caKeyFile := path.Join(o.OutputDir, o.CAName+cert.KeyFileSuffix)
	caCertFile := path.Join(o.OutputDir, o.CAName+cert.CertFileSuffix)
	if _, err := os.Stat(caKeyFile); !errors.Is(err, os.ErrNotExist) {
		msg := fmt.Sprintf("ca key file (%s) already exist, remove it if you want to continue", caKeyFile)
		return errors.New(msg)
	}
	if _, err := os.Stat(caCertFile); !errors.Is(err, os.ErrNotExist) {
		msg := fmt.Sprintf("ca cert file (%s) already exist, remove it if you want to continue", caCertFile)
		return errors.New(msg)
	}

	fmt.Printf("generating ca key and cert: %s\n", o.CAName)
	kc, err := cert.NewCA(o.CAName, CertValidDays)
	if err != nil {
		msg := fmt.Sprintf("generate CA key and cert failed, err: %s", err)
		return errors.New(msg)
	}
	kc.Dump(o.OutputDir)
	return nil
}

func (o *GenCertOptions) genCert() error {
	caKeyFile := path.Join(o.OutputDir, o.CAName+cert.KeyFileSuffix)
	caCertFile := path.Join(o.OutputDir, o.CAName+cert.CertFileSuffix)
	caKC, err := cert.LoadKeyCertPEMFile(caKeyFile, caCertFile)
	if err != nil {
		msg := fmt.Sprintf("load CA key and cert file failed, err: %s", err)
		return errors.New(msg)
	}
	for _, name := range o.Names {
		fmt.Printf("generating key and cert for: %s\n", name)

		kc := cert.NewKeyCert(name)
		if err := kc.GenKey(); err != nil {
			msg := fmt.Sprintf("generate private key for %s failed, err: %s", name, err)
			return errors.New(msg)
		}
		if err := kc.SignedByCA(caKC, CertValidDays); err != nil {
			msg := fmt.Sprintf("the CA sign the cert failed, err: %s", err)
			return errors.New(msg)
		}
		kc.Dump(o.OutputDir)
	}
	return nil
}

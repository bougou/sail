package x

import (
	"fmt"

	"github.com/bougou/gopkg/common"
	"github.com/bougou/sail/pkg/commands/x/gencert"
	"github.com/bougou/sail/pkg/models"
	"github.com/spf13/cobra"
)

func NewCmdX(sailOption *models.SailOption) *cobra.Command {
	o := NewXOptions(sailOption)

	cmd := &cobra.Command{
		Use:   "x",
		Short: "helper tools",
		Long:  "helper tools",
		Run: func(cmd *cobra.Command, args []string) {
			common.CheckErr(o.Complete(cmd, args))
			common.CheckErr(o.Validate())
			common.CheckErr(o.Run(args))
		},
	}

	cmd.AddCommand(gencert.NewCmdGenCert(o.sailOption))

	return cmd
}

type XOptions struct {
	sailOption *models.SailOption
}

func NewXOptions(sailOption *models.SailOption) *XOptions {
	return &XOptions{
		sailOption: sailOption,
	}
}

func (o *XOptions) Complete(cmd *cobra.Command, args []string) error {
	return nil
}

func (o *XOptions) Validate() error {
	return nil
}

func (o *XOptions) Run(args []string) error {
	fmt.Println("Specify a concret command under x")
	return nil
}

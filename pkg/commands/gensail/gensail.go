package gensail

import (
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/bougou/gopkg/common"
	"github.com/bougou/sail/pkg/models"
	"github.com/bougou/sail/pkg/models/product"
	"github.com/spf13/cobra"
)

func NewCmdGenSail(sailOption *models.SailOption) *cobra.Command {
	o := NewGenSailOptions(sailOption)

	cmd := &cobra.Command{
		Use:   "gen-sail",
		Short: "auto generate the sail.yaml playbook file",
		Long:  "auto generate the sail.yaml playbook file",
		Run: func(cmd *cobra.Command, args []string) {
			common.CheckErr(o.Complete(cmd, args))
			common.CheckErr(o.Validate())
			common.CheckErr(o.Run())
		},
	}

	defaultProductName := ""
	cmd.Flags().StringVarP(&o.productName, "product", "p", defaultProductName, "the product name")
	_ = cmd.MarkFlagRequired("playbook")

	return cmd
}

type GenSailOptions struct {
	productName string
	productDir  string

	sailOption *models.SailOption
}

func NewGenSailOptions(sailOption *models.SailOption) *GenSailOptions {
	return &GenSailOptions{
		sailOption: sailOption,
	}
}

func (o *GenSailOptions) Complete(cmd *cobra.Command, args []string) error {
	if o.productName == "" {
		return errors.New("product name must not be empty")
	}

	o.productDir = path.Join(o.sailOption.ProductsDir, o.productName)
	stat, err := os.Stat(o.productDir)
	if err != nil || !stat.IsDir() {
		return fmt.Errorf("not found dir of product, %s does not exist", o.productDir)
	}

	return nil
}

func (o *GenSailOptions) Validate() error {
	return nil
}

func (o *GenSailOptions) Run() error {
	product := product.NewProduct(o.productName, o.sailOption.ProductsDir)
	if err := product.Init(); err != nil {
		return fmt.Errorf("product init failed, err: %s", err)
	}

	playbook, err := product.GenSail()
	if err != nil {
		return fmt.Errorf("gen sail playbook failed, err: %s", err)
	}

	b, err := common.Encode("yaml", playbook)
	if err != nil {
		fmt.Println("encode vars failed", err)
	}

	if err := os.WriteFile(product.SailPlaybookFile(), b, 0644); err != nil {
		fmt.Println("write product sail playbook file failed", err)
	}

	return nil
}

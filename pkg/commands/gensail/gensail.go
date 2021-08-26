package gensail

import (
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/bougou/gopkg/common"
	"github.com/bougou/sail/pkg/models"
	"github.com/spf13/cobra"
)

func NewCmdGenSail(sailOption *models.SailOption) *cobra.Command {
	o := NewGenSailOptions(sailOption)

	cmd := &cobra.Command{
		Use:   "gen-sail",
		Short: "gen-sail",
		Long:  "gen-sail",
		Run: func(cmd *cobra.Command, args []string) {
			common.CheckErr(o.Complete(cmd, args))
			common.CheckErr(o.Validate())
			common.CheckErr(o.Run())
		},
	}

	defaultProductName := ""
	cmd.Flags().StringVarP(&o.productName, "product", "p", defaultProductName, "the product name")
	cmd.MarkFlagRequired("playbook")

	return cmd
}

type GenSailOptions struct {
	productName string
	productDir  string

	args        []string
	productsDir string
	targetsDir  string
	sailOption  *models.SailOption
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
		msg := fmt.Sprintf("Not found dir of product, %s does not exist", o.productDir)
		return errors.New(msg)
	}

	return nil
}

func (o *GenSailOptions) Validate() error {
	return nil
}

func (o *GenSailOptions) Run() error {
	product := models.NewProduct(o.productName, o.sailOption.ProductsDir)
	if err := product.Init(); err != nil {
		msg := fmt.Sprintf("product init failed, err: %s", err)
		return errors.New(msg)
	}

	playbook, err := product.GenSail()
	if err != nil {
		msg := fmt.Sprintf("GenSail failed, err: %s", err)
		return errors.New(msg)
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

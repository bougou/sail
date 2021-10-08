package bundle

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

func NewCmdBundle(sailOption *models.SailOption) *cobra.Command {
	o := NewBundleOptions(sailOption)

	cmd := &cobra.Command{
		Use:   "bundle",
		Short: "bundle packages of a product",
		Long:  "bundle packages of a product",
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

type BundleOptions struct {
	productName string
	productDir  string

	args        []string
	productsDir string
	targetsDir  string
	sailOption  *models.SailOption
}

func NewBundleOptions(sailOption *models.SailOption) *BundleOptions {
	return &BundleOptions{
		sailOption: sailOption,
	}
}

func (o *BundleOptions) Complete(cmd *cobra.Command, args []string) error {
	if o.productName == "" {
		return errors.New("product name must not be empty")
	}

	o.productDir = path.Join(o.sailOption.ProductsDir, o.productName)
	stat, err := os.Stat(o.productDir)
	if err != nil || !stat.IsDir() {
		msg := fmt.Sprintf("not found dir of product, %s does not exist", o.productDir)
		return errors.New(msg)
	}

	return nil
}

func (o *BundleOptions) Validate() error {
	return nil
}

func (o *BundleOptions) Run() error {
	product := product.NewProduct(o.productName, o.sailOption.ProductsDir)
	if err := product.Init(); err != nil {
		msg := fmt.Sprintf("product init failed, err: %s", err)
		return errors.New(msg)
	}

	components := product.ComponentList()
	fmt.Printf("the product %s contains (%d) components:\n", o.productName, len(components))
	for _, c := range components {
		fmt.Printf("- %s\n", c)
	}

	return nil
}

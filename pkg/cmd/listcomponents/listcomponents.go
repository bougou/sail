package listcomponents

import (
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/bougou/sail/pkg/models"
	cmdutil "github.com/bougou/sail/pkg/util"
	"github.com/spf13/cobra"
)

func NewCmdListComponents(sailOption *models.SailOption) *cobra.Command {
	o := NewListComponentsOptions(sailOption)

	cmd := &cobra.Command{
		Use:   "list-components",
		Short: "list-components",
		Long:  "list-components",
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(o.Complete(cmd, args))
			cmdutil.CheckErr(o.Run())
		},
	}

	defaultProductName := ""
	cmd.Flags().StringVarP(&o.productName, "product", "p", defaultProductName, "the product name")
	cmd.MarkFlagRequired("playbook")

	return cmd
}

type ListComponentsOptions struct {
	productName string
	productDir  string

	args        []string
	productsDir string
	targetsDir  string
	sailOption  *models.SailOption
}

func NewListComponentsOptions(sailOption *models.SailOption) *ListComponentsOptions {
	return &ListComponentsOptions{
		sailOption: sailOption,
	}
}

func (o *ListComponentsOptions) Complete(cmd *cobra.Command, args []string) error {
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

func (o *ListComponentsOptions) Run() error {
	product := models.NewProduct(o.productName, o.sailOption.ProductsDir)
	if err := product.Init(); err != nil {
		msg := fmt.Sprintf("product init failed, err: %s", err)
		return errors.New(msg)
	}

	components := product.ComponentList()
	fmt.Printf("The product %s contains the following components:\n", o.productName)
	for _, c := range components {
		fmt.Printf("- %s\n", c)
	}

	return nil
}

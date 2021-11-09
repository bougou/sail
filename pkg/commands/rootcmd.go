package commands

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/bougou/sail/pkg/commands/apply"
	"github.com/bougou/sail/pkg/commands/confcreate"
	"github.com/bougou/sail/pkg/commands/confupdate"
	"github.com/bougou/sail/pkg/commands/gensail"
	"github.com/bougou/sail/pkg/commands/listcomponents"
	"github.com/bougou/sail/pkg/commands/upgrade"
	"github.com/bougou/sail/pkg/commands/x"
	"github.com/bougou/sail/pkg/models"
	"github.com/bougou/sail/pkg/version"
	"github.com/mitchellh/go-homedir"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const defaultSailRC string = ".sailrc"
const envPrefix = "SAIL"
const homePage = "https://github.com/bougou/sail"

var (
	cfgFile string
)

func NewSailCommand() *cobra.Command {
	sailOption := &models.SailOption{}
	showVersion := false

	rootCmd := &cobra.Command{
		Use:   "sail",
		Short: "sail",
		Long:  fmt.Sprintf("sail\n\nFind more information at: %s\n\n", homePage),
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if showVersion {
				fmt.Printf("Version: %s\n", version.Version)
				fmt.Printf("Commit: %s\n", version.Commit)
				fmt.Printf("BuildAt: %s\n", version.BuildAt)
				return nil
			}

			return initConfig(cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {
			if showVersion {
				return
			}

			fmt.Println("you have to specify a subcommand for sail")
		},
	}

	execFile, err := os.Executable()
	if err != nil {
		panic("fatal")
	}

	defaultTargetsDir := filepath.Join(filepath.Dir(execFile), "targets")
	defaultProductsDir := filepath.Join(filepath.Dir(execFile), "products")
	defaultPackagesDir := filepath.Join(filepath.Dir(execFile), "packages")

	rootCmd.PersistentFlags().BoolVarP(&showVersion, "version", "v", false, "show version")
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "", "", "config file, default $HOME/.sailrc")
	rootCmd.PersistentFlags().StringVarP(&sailOption.TargetsDir, "targets-dir", "", defaultTargetsDir, "the targets dir")
	rootCmd.PersistentFlags().StringVarP(&sailOption.ProductsDir, "products-dir", "", defaultProductsDir, "the products dir")
	rootCmd.PersistentFlags().StringVarP(&sailOption.PackagesDir, "packages-dir", "", defaultPackagesDir, "the packages dir")

	rootCmd.PersistentFlags().StringVarP(&sailOption.DefaultTarget, "default-target", "", "", "the default target")
	rootCmd.PersistentFlags().StringVarP(&sailOption.DefaultZone, "default-zone", "", "", "the default zone")

	rootCmd.Flags().AddGoFlagSet(flag.CommandLine)

	rootCmd.AddCommand(apply.NewCmdApply(sailOption))
	rootCmd.AddCommand(confcreate.NewCmdConfCreate(sailOption))
	rootCmd.AddCommand(confupdate.NewCmdConfUpdate(sailOption))
	rootCmd.AddCommand(gensail.NewCmdGenSail(sailOption))
	rootCmd.AddCommand(listcomponents.NewCmdListComponents(sailOption))
	rootCmd.AddCommand(upgrade.NewCmdUpgrade(sailOption))
	rootCmd.AddCommand(x.NewCmdX(sailOption))

	return rootCmd
}

func initConfig(cmd *cobra.Command) error {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		viper.AddConfigPath(home)
		viper.SetConfigName(defaultSailRC)
	}

	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// fmt.Println("not found any config file")
		} else {
			// ConfigFile found, but viper read failed.
			fmt.Println("using config file:", viper.ConfigFileUsed())
			fmt.Println(err)
		}
	} else {
		fmt.Println("using config file:", viper.ConfigFileUsed())
	}

	// When we bind flags to environment variables expect that the
	// environment variables are prefixed, e.g. a flag like --targets-dir
	// binds to an environment variable SAIL_TARGETS_DIR. This helps
	// avoid conflicts.
	viper.SetEnvPrefix(envPrefix)

	// Bind to environment variables
	// Works great for simple config names, but needs help for names contains hyphen
	// like --favorite-color which we fix in the bindFlags function
	viper.AutomaticEnv()

	// Bind the current command's flags to viper
	bindFlags(cmd, viper.GetViper())

	return nil
}

// Use viper's value if found to SET cobra flags which are not set(NOT CHANGED).
func bindFlags(cmd *cobra.Command, v *viper.Viper) {
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		// Environment variables can't have dashes in them, so bind them to their equivalent
		// keys with underscores, e.g. --favorite-color to <envPrefix>_FAVORITE_COLOR
		if strings.Contains(f.Name, "-") {
			envVarSuffix := strings.ToUpper(strings.ReplaceAll(f.Name, "-", "_"))
			v.BindEnv(f.Name, fmt.Sprintf("%s_%s", envPrefix, envVarSuffix))
		}

		// Apply the viper config value to the flag when the flag is not set and viper has a value
		// If f.Changed is true which means the the flag is specified when run the cmd,
		// then so just use it, because its priority is highest.
		if !f.Changed && v.IsSet(f.Name) {
			val := v.Get(f.Name)
			cmd.Flags().Set(f.Name, fmt.Sprintf("%v", val))
		}
	})
}

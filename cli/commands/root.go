package commands

import (
	"errors"
	"fmt"
	"os"

	"github.com/myerscode/aws-meta/cli/internal/util"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "aws-meta",
	Short: "A tool for collecting and looking at information about AWS Partitions and Regions.",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		envVarBind(cmd)
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.aws-meta.yaml)")
	rootCmd.Flags().BoolP("trace", "t", false, "Show full trace logs.")
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".aws-meta")
	}

	viper.SetEnvPrefix("AWSMETA")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		// A missing config file is fine; the CLI works without one. Only
		// surface genuine read failures (bad YAML, permissions, or an
		// explicitly requested file that cannot be read).
		var notFound viper.ConfigFileNotFoundError
		if !errors.As(err, &notFound) {
			util.LogError(fmt.Sprintf("Error using config file %s: %s", viper.ConfigFileUsed(), err))
		}
	}
}

func envVarBind(cmd *cobra.Command) {
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		configName := f.Name
		if !f.Changed && viper.IsSet(configName) {
			val := viper.Get(configName)
			err := cmd.Flags().Set(f.Name, fmt.Sprintf("%v", val))
			if err != nil {
				util.LogError(err.Error())
			}
		}
	})
}

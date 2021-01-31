package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "ndh",
	Short: "Get hierarchical dependencies for node packages",
	Long: "Given a node package name show hierarchical/tree dependencies of the package",
	RunE: func(_ *cobra.Command, args []string) error {
		return nil
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	//TODO: Implement loading .env config
}

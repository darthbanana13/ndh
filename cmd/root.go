package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

var cfgFile string

type npmPackage struct {
	Name			string				`json:"name"`
	Dependencies	map[string]string	`json:"dependencies"`
}

var rootCmd = &cobra.Command{
	Use:   "ndh",
	Short: "Get hierarchical dependencies for node packages",
	Long: "Given a node package name show hierarchical/tree dependencies of the package",
	RunE: func(_ *cobra.Command, args []string) error {
		pkgName := args[0]
		pkgVer := args[1]
		//TODO: Handle errors
		resp, _ := http.Get("https://registry.npmjs.org/" + pkgName + "/" + pkgVer)
		body, _ := ioutil.ReadAll(resp.Body)
		var pak npmPackage
		_ = json.Unmarshal(body, &pak)
		fmt.Println(pak)
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

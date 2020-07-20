package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"lingo/utils"
	"os"
)

var (
	serverURIFlag,
	serverIDFlag,
	sourceServerFlag,
	targetServerFlag string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "lingo",
	Short: "",
	Long:  ``,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(utils.ColorfulString("red", err.Error()))
		os.Exit(1)
	}
}

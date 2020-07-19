package cmd

import (
	"errors"
	"github.com/spf13/cobra"
	"mongo-transfer/utils"
	"os"
)

var (
	serverURIFlag,
	serverIDFlag,
	sourceServerFlag,
	targetServerFlag string
)

var ErrRecordNotFound = errors.New("Unable to find file or directory named servers/local")

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "mongo-transfer",
	Short: "",
	Long:  ``,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		utils.ColorfulString("red", err.Error())
		os.Exit(1)
	}
}

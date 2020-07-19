package cmd

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"mongo-transfer/database"
	"mongo-transfer/transfer"
	"mongo-transfer/utils"
	"time"
)

// transferCmd represents the transfer command
var transferCmd = &cobra.Command{
	Use:   "transfer",
	Short: "",
	Long:  ``,
	Args: func(cmd *cobra.Command, args []string) error {
		if sourceServerFlag == "" && targetServerFlag == "" {
			return errors.New("both from and to flags must be filled")
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		fmt.Printf("Starting transfer from %s to %s\n", sourceServerFlag, targetServerFlag)

		sourceServer := database.Server{ID: sourceServerFlag}
		targetServer := database.Server{ID: targetServerFlag}

		if err = utils.StepsFunctions(
			sourceServer.Fetch,
			sourceServer.Ping,
			targetServer.Fetch,
			targetServer.Ping); err != nil {
			return
		}

		fmt.Println(utils.ColorfulString("green", "Successfully connected to both servers"))

		transfer := transfer.Transfer{
			Source:    &sourceServer,
			Target:    &targetServer,
			StartedAt: time.Now(),
		}

		finishedAt, err := transfer.Start()
		if err != nil {
			return
		}

		fmt.Println(utils.ColorfulString("green", "\nOperation finished"))

		elapsed := finishedAt.Sub(transfer.StartedAt)
		fmt.Println("Took", elapsed)

		return
	},
}

func init() {
	rootCmd.AddCommand(transferCmd)

	transferCmd.Flags().StringVar(&sourceServerFlag, "from", "", "ID of the source server")
	transferCmd.Flags().StringVar(&targetServerFlag, "to", "", "ID of the target server")
}

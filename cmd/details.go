package cmd

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"mongo-transfer/database"
)

var detailsCmd = &cobra.Command{
	Use:   "details",
	Short: "",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		if args[0] == "" {
			return errors.New("not a valid server ID")
		}

		serverID := args[0]

		fmt.Printf("Fetching details for \"%s\"\n", serverID)

		server := database.Server{}

		err = server.Load(serverID)
		if err != nil {
			return
		}

		defer server.Disconnect()

		err = server.LoadAll()
		if err != nil {
			return
		}

		fmt.Println(server.Databases)

		return
	},
}

func init() {
	rootCmd.AddCommand(detailsCmd)
}

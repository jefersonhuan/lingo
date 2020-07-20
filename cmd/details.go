package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"lingo/database"
)

var detailsCmd = &cobra.Command{
	Use:   "details",
	Short: "",
	Long:  ``,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		serverID := args[0]

		fmt.Printf("Fetching details for \"%s\"\n", serverID)

		server := database.Server{ID: serverID}

		err = server.Fetch()
		if err != nil {
			return
		}

		defer server.Disconnect()

		err = server.LoadAll()
		if err != nil {
			return
		}

		if len(server.Databases) == 0 {
			fmt.Println("No Databases where found")
		}

		for _, database := range server.Databases {
			fmt.Println(database)
		}

		return
	},
}

func init() {
	rootCmd.AddCommand(detailsCmd)
}

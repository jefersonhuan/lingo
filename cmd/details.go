package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"mongo-transfer/database"
)

var listDatabasesCmd = &cobra.Command{
	Use:   "details",
	Short: "",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		server := database.Server{}

		err = server.Load("local")
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
	rootCmd.AddCommand(listDatabasesCmd)
}

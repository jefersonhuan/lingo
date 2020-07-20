package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"lingo/database"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "delete a previously saved MongoDB server",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		server := database.Server{ID: args[0]}

		err := server.Delete()
		if err != nil {
			return err
		}

		fmt.Println("Successfully deleted", server.ID)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}

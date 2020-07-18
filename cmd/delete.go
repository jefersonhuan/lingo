package cmd

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"mongo-transfer/database"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "delete a previously saved MongoDB server",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		server := database.Server{ID: args[0]}

		err := server.Delete()
		if err == ErrRecordNotFound {
			return errors.New("server not found")
		} else if err != nil {
			return err
		}

		fmt.Println("Successfully deleted", server.ID)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}

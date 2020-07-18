package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"mongo-transfer/database"
	"time"
)

var addServerCmd = &cobra.Command{
	Use:   "addServer",
	Short: "Add a mongodb Server",
	Long: `Ex. 
mongo-transfer addServer --from-uri="mongodb://localhost:27017"`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fn := func() func() error {
			if len(args) == 0 {
				return AddServerPrompt
			}

			return nil
		}()

		if err := fn(); err != nil {
			return err
		}

		return nil
	},
}

func AddServerPrompt() (err error) {
	var host, username, password string
	var port int

	var argsRead int

	fmt.Print("Host [default localhost]: ")
	if argsRead, err = fmt.Scanln(&host); argsRead == 0 {
		host = "localhost"
	} else if err != nil {
		return
	}

	fmt.Print("Port [default 27017]: ")
	if argsRead, err = fmt.Scanln(&port); argsRead == 0 {
		port = 27017
	} else if err != nil {
		return
	}

	fmt.Print("User [default blank]: ")
	if argsRead, err = fmt.Scanln(&username); argsRead != 0 && err != nil {
		return
	}

	fmt.Print("Password [default blank]: ")
	if argsRead, err = fmt.Scanln(&password); argsRead != 0 && err != nil {
		return
	}

	fmt.Println("Now, we need an identifier to this server")

	server := database.Server{}
	server.BuildURI(username, password, host, port)

	for {
		fmt.Print("ID: ")
		if argsRead, err = fmt.Scanln(&server.ID); argsRead == 0 || err != nil {
			fmt.Println("Please, insert a valid ID")
			continue
		}

		break
	}

	server.CreatedAt = time.Now()

	var option string
	fmt.Print("Do you wish to test connection now? [y/n] ")
	if _, err = fmt.Scanf("%s", &option); err != nil {
		fmt.Println("Oops! Couldn't test the server, but we'll save either way")
	} else if option == "y" {
		fmt.Println("Testing connection...")
		if err = server.Ping(); err != nil {
			fmt.Println("Couldn't connect to the server, because:", err.Error())
			fmt.Println("We'll save either way. But you can change (or delete) this server at anytime")
		} else {
			fmt.Printf("Successfully connected to %s\n", server.ID)
			server.Disconnect()
		}
	}

	if err = server.Save(); err != nil {
		return
	}

	fmt.Println("Successfully saved server")

	return
}

func init() {
	rootCmd.AddCommand(addServerCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addServerCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addServerCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

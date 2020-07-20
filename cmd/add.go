package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"lingo/database"
	"lingo/utils"
	"time"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a mongodb Server",
	Long: `Ex. 
lingo add --from-uri="mongodb://localhost:27017" --name="local"`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fn := func() func() error {
			if serverURIFlag != "" {
				return AddServerFromURI
			}

			return AddServerPrompt
		}()

		if err := fn(); err != nil {
			return err
		}

		return nil
	},
}

func AddServerFromURI() error {
	server := &database.Server{URI: serverURIFlag, ID: serverIDFlag}

	return testAndSave(server)
}

func AddServerPrompt() (err error) {
	var host, username, password string
	var port int

	var argsRead int

	fmt.Print("Host: [default localhost] ")
	if argsRead, err = fmt.Scanln(&host); argsRead == 0 {
		host = "localhost"
	} else if err != nil {
		return
	}

	fmt.Print("Port: [default 27017] ")
	if argsRead, err = fmt.Scanln(&port); argsRead == 0 {
		port = 27017
	} else if err != nil {
		return
	}

	fmt.Print("User: [default blank] ")
	if argsRead, err = fmt.Scanln(&username); argsRead != 0 && err != nil {
		return
	}

	fmt.Print("Password: [default blank] ")
	if argsRead, err = fmt.Scanln(&password); argsRead != 0 && err != nil {
		return
	}

	server := database.Server{}
	server.BuildURI(username, password, host, port)

	return testAndSave(&server)
}

func testAndSave(server *database.Server) (err error) {
	var argsRead int

	if server.ID == "" {
		fmt.Println("Please, insert an identifier for this server")

		for {
			fmt.Print("ID: ")
			if argsRead, err = fmt.Scanln(&server.ID); argsRead == 0 || err != nil {
				fmt.Println("Please, insert a valid ID")
				continue
			}

			break
		}
	}

	server.CreatedAt = time.Now()

	if err = testConnectionPrompt(server); err != nil {
		return
	}

	if err = server.Save(); err != nil {
		return
	}

	fmt.Println("Successfully saved server")

	return
}

func testConnectionPrompt(server *database.Server) (err error) {
	var option string

	fmt.Print("Do you want to test the connection right now? [y/n] ")

	if _, err = fmt.Scanf("%s", &option); err != nil {
		fmt.Println("Oops! Couldn't test the server, but we'll save it either way")
	} else if option == "y" {
		fmt.Println("Testing connection...")

		if err = server.Ping(); err != nil {
			fmt.Println("Couldn't connect to the server, because:", err.Error())
			fmt.Println("We'll save either way. But you can change (or delete) this server at anytime")
		} else {
			fmt.Println(utils.ColorfulString("green", "Successfully connected to "+server.URI))
		}
	}

	return
}

func init() {
	rootCmd.AddCommand(addCmd)

	addCmd.Flags().StringVar(&serverURIFlag, "from-uri", "", "adds server with given URI parameter")
	addCmd.Flags().StringVar(&serverIDFlag, "name", "", "flag for the server's name")
}

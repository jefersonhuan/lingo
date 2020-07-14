/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"mongo-transfer/models"
	"mongo-transfer/operations"
	"time"
)

var addServerCmd = &cobra.Command{
	Use:   "addServer",
	Short: "Add a mongodb Server",
	Long: `Ex. 
mongo-transfer addServer --from-uri="mongodb://..."`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fn := func() func() error {
			if len(args) == 0 {
				return addServerPrompt
			}

			return nil
		}()

		if err := fn(); err != nil {
			return err
		}

		return nil
	},
}

func addServerPrompt() (err error) {
	server := models.Server{}

	var argsRead int

	fmt.Print("Host [default localhost]: ")
	if argsRead, err = fmt.Scanln(&server.Host); argsRead == 0 {
		server.Host = "localhost"
	} else if err != nil {
		return
	}

	fmt.Print("Port [default 27017]: ")
	if argsRead, err = fmt.Scanln(&server.Port); argsRead == 0 {
		server.Port = 27017
	} else if err != nil {
		return
	}

	fmt.Print("User [default blank]: ")
	if argsRead, err = fmt.Scanln(&server.User); argsRead != 0 && err != nil {
		return
	}

	fmt.Print("Password [default blank]: ")
	if argsRead, err = fmt.Scanln(&server.Password); argsRead != 0 && err != nil {
		return
	}

	fmt.Println("Now, we need an identifier to this server")

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
		if err = operations.TestConnection(server); err != nil {
			fmt.Println("Couldn't connect to the server, because:", err.Error())
			fmt.Println("We'll save either way. But you can change (or delete) this server at anytime")
		} else {
			fmt.Printf("Successfully connected to %s\n", server.ID)
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

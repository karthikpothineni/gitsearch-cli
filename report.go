package main

import (
	"fmt"
	"gitsearch-cli/handlers"
	"gopkg.in/ukautz/clif.v1"
)

func callBackFunction(c *clif.Command) {
	org := c.Option("organization").String()
	authKey := c.Option("auth-key").String()

	fmt.Println("hi")
	fmt.Println(org)
	fmt.Println(authKey)

	// create git handler
	gitHandler := handlers.NewRequestHandler(authKey)

	// list repositories
	gitHandler.ListRepositories()
}

func main() {
	// add commands and options
	cli := clif.New("report app", "1.0.0", "Fetches git information based on organisation name")
	cmd := clif.NewCommand("report", "Pass organization and auth-key", callBackFunction).
		NewOption("organization", "o", "Name of the organization", "", true, false).
		NewOption("auth-key", "a", "Github Auth Key", "", true, false)
	cli.Add(cmd)
	cli.Run()


}

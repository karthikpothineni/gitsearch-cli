package main

import (
	"fmt"
	"gitsearch-cli/handlers"
	"gopkg.in/ukautz/clif.v1"
)

// report constants
const (
	defaultProgressBarSize        = 100
	defaultProgressBarRenderWidth = 80
)

// callBackFunction - callback function for command
func callBackFunction(c *clif.Command) {
	orgName := c.Option("organization").String()
	authKey := c.Option("auth-key").String()

	// create git handler
	gitHandler := handlers.NewRequestHandler(authKey)

	// start progress bar
	fmt.Println("Progress:")
	progressBar := clif.NewProgressBar(defaultProgressBarSize).SetStyle(clif.ProgressBarStyleAscii)
	progressBar.SetRenderWidth(defaultProgressBarRenderWidth)

	// get repo information from github
	contributorRepoDetails, repoLanguageDetails, err := gitHandler.FetchRepositoriesInfo(orgName, progressBar)
	if err != nil {
		fmt.Println(err)
		return
	}

	// finish progress bar
	progressBar.Finish()

	// handle response
	gitHandler.HandleResponse(contributorRepoDetails, repoLanguageDetails)
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

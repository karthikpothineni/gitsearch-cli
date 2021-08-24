package main

import (
	"errors"
	"fmt"
	"strings"

	"gopkg.in/ukautz/clif.v1"

	"gitsearch-cli/handlers"
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

	// validate options
	if err := validateOptions(orgName, authKey); err != nil {
		fmt.Println(err)
		return
	}

	// create git handler
	gitHandler := handlers.NewRequestHandler(authKey)

	// start progress bar
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

// validateOptions - validates the options
func validateOptions(orgName, authKey string) error {
	if strings.TrimSpace(orgName) == "" {
		return errors.New("error: organization name cannot be empty")
	}

	if strings.TrimSpace(authKey) == "" {
		return errors.New("error: auth key cannot be empty")
	}
	return nil
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

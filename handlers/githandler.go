package handlers

import (
	"context"
	"errors"
	"fmt"
	"gitsearch-cli/utils"
	"gopkg.in/ukautz/clif.v1"
	"strings"

	"github.com/google/go-github/v38/github"
	"golang.org/x/oauth2"
)

// repo constants
const joinSeparator = ", "

// RequestHandler - holds git client
type RequestHandler struct {
	GitClient   *github.Client
	UserContext context.Context
}

// NewRequestHandler - returns a new request handler object
func NewRequestHandler(accessToken string) (*RequestHandler, error) {
	if strings.TrimSpace(accessToken) == "" {
		return nil, errors.New("invalid access token")
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accessToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	return &RequestHandler{
		GitClient:   client,
		UserContext: ctx,
	}, nil
}

// FetchRepositoriesInfo - returns all the repository information including contributors and languages
func (r *RequestHandler) FetchRepositoriesInfo(orgName string, progressBar clif.ProgressBar) (map[string][]string, map[string][]string, error) {
	// list all repositories for the authenticated user
	repos, _, err := r.GitClient.Repositories.ListByOrg(r.UserContext, orgName, nil)
	if err != nil {
		if _, ok := err.(*github.RateLimitError); ok {
			err = errors.New("hit rate limit. only 5000 requests per hour are allowed")
			return nil, nil, err
		}
		err = fmt.Errorf("error occurred while listing repos: %s", err)
		return nil, nil, err
	}

	// change progress bar size
	progressBar.SetSize(len(repos))

	// get contributors and language info
	contributorRepoDetails, repoLanguageDetails, err := r.fetchRepoContributorsAndLanguages(repos, progressBar, orgName)

	return contributorRepoDetails, repoLanguageDetails, err
}

// fetchRepoContributorsAndLanguages - returns contributor and language information for each repo
func (r *RequestHandler) fetchRepoContributorsAndLanguages(repos []*github.Repository, progressBar clif.ProgressBar, orgName string) (map[string][]string, map[string][]string, error) {
	contributorRepoDetails := make(map[string][]string)
	repoLanguageDetails := make(map[string][]string)
	var repoErr error

	// get contributors and language info
	fmt.Println("Progress:")
	for _, eachRepo := range repos {
		fmt.Print("\r" + progressBar.Render())
		_ = progressBar.Increment()
		if eachRepo.Name != nil {

			// get contributors for a repo
			contributors, _, err := r.GitClient.Repositories.ListContributors(r.UserContext, orgName, *eachRepo.Name, nil)
			if err == nil {
				for _, eachContributor := range contributors {
					if eachContributor.Login != nil {
						var contributorDetails string
						userData, _, err := r.GitClient.Users.Get(r.UserContext, *eachContributor.Login)
						if err != nil {
							contributorDetails = fmt.Sprintf("%s;%s;%s", *eachContributor.Login, "", "")
						} else {
							contributorDetails = fmt.Sprintf("%s;%s;%s", *eachContributor.Login, utils.HandleNilString(userData.Name), utils.HandleNilString(userData.Email))
						}
						if repos, ok := contributorRepoDetails[contributorDetails]; ok {
							repos = append(repos, *eachRepo.Name)
							contributorRepoDetails[contributorDetails] = repos
						} else {
							contributorRepoDetails[contributorDetails] = []string{*eachRepo.Name}
						}
					}
				}
			} else {
				repoErr = fmt.Errorf("unable to list contributors: %s", err.Error())
			}

			// get languages for a repo
			languages, _, err := r.GitClient.Repositories.ListLanguages(r.UserContext, orgName, *eachRepo.Name)
			if err == nil {
				repoLanguageDetails[*eachRepo.Name] = utils.GetMapKeys(languages)
			} else {
				repoErr = fmt.Errorf("unable to list languages: %s", err.Error())
			}
		}
	}
	fmt.Print("\r" + progressBar.Render() + "\n")

	return contributorRepoDetails, repoLanguageDetails, repoErr
}

// HandleResponse - prepare and print the response
func (r *RequestHandler) HandleResponse(contributorRepoDetails map[string][]string, repoLanguageDetails map[string][]string) {
	for eachContributor, repos := range contributorRepoDetails {
		var repoDetails string
		var allLanguages []string
		for _, eachRepo := range repos {
			repoDetails += eachRepo + joinSeparator
			if languages, ok := repoLanguageDetails[eachRepo]; ok {
				allLanguages = append(allLanguages, languages...)
			}
		}
		repoDetails = strings.TrimSuffix(repoDetails, joinSeparator)
		allLanguages = utils.RemoveDuplicates(allLanguages)
		languageDetails := strings.Join(allLanguages[:], joinSeparator)
		response := eachContributor + ";" + repoDetails + ";" + languageDetails
		fmt.Printf(response + "\n")
	}
}

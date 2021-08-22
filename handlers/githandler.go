package handlers

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/google/go-github/v38/github"
	"golang.org/x/oauth2"
)

// RequestHandler - holds git client
type RequestHandler struct {
	GitClient   *github.Client
	UserContext context.Context
}

// NewRequestHandler - returns a new request handler object
func NewRequestHandler(accessToken string) *RequestHandler {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accessToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	return &RequestHandler{
		GitClient:   client,
		UserContext: ctx,
	}
}

func (r *RequestHandler) ListRepositories(orgName string) (map[string][]string, map[string][]string, error) {
	contributorRepoDetails := make(map[string][]string)
	repoLanguageDetails := make(map[string][]string)

	// list all repositories for the authenticated user
	repos, _, err := r.GitClient.Repositories.ListByOrg(r.UserContext, orgName, nil)
	if err != nil {
		if _, ok := err.(*github.RateLimitError); ok {
			err = errors.New("Hit rate limit. Only 5000 requests per hour are allowed")
			return nil, nil, err
		}
		err = fmt.Errorf("Error occurred while listing repos: %s", err.Error())
		return nil, nil, err
	}
	for _, eachRepo := range repos {
		if eachRepo.Name != nil {
			// get contributors for a repo
			contributors, _, err := r.GitClient.Repositories.ListContributors(r.UserContext, orgName, *eachRepo.Name, nil)
			if err == nil {
				for _, eachContributor := range contributors {
					contributorDetails := fmt.Sprintf("%s;%s;%s", handleNilString(eachContributor.Login), handleNilString(eachContributor.Name), handleNilString(eachContributor.Email))
					if repos, ok := contributorRepoDetails[contributorDetails]; ok {
						repos = append(repos, *eachRepo.Name)
						contributorRepoDetails[contributorDetails] = repos
					} else {
						contributorRepoDetails[contributorDetails] = []string{*eachRepo.Name}
					}
				}
			}
			// get languages for a repo
			languages, _, err := r.GitClient.Repositories.ListLanguages(r.UserContext, orgName, *eachRepo.Name)
			if err == nil {
				repoLanguageDetails[*eachRepo.Name] = GetMapKeys(languages)
			}
		}
	}

	return contributorRepoDetails, repoLanguageDetails, err
}

// HandleResponse - prepare and print the response
func (r *RequestHandler) HandleResponse(contributorRepoDetails map[string][]string, repoLanguageDetails map[string][]string) {

	for eachContributor, repos := range contributorRepoDetails {
		var repoDetails, languageDetails string
		for _, eachRepo := range repos {
			repoDetails += eachRepo + ", "
			if languages, ok := repoLanguageDetails[eachRepo]; ok {
				languageDetails = strings.Join(languages[:], ", ")
			}
		}
		response := eachContributor + ";" + repoDetails + ";" + languageDetails
		fmt.Printf(response + "\n")
	}
}

// GetMapKeys - returns keys in a map
func GetMapKeys(v interface{}) []string {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Map {
		return nil
	}
	t := rv.Type()
	if t.Key().Kind() != reflect.String {
		return nil
	}
	var result []string
	for _, kv := range rv.MapKeys() {
		result = append(result, kv.String())
	}
	return result
}

// handleNilString - return empty string if string pointer is nil
func handleNilString(value *string) string {
	if value == nil {
		return ""
	} else {
		return *value
	}
}

package handlers

import (
	"context"
	"fmt"

	"golang.org/x/oauth2"
	"github.com/google/go-github/v38/github"
)

// RequestHandler - holds git client
type RequestHandler struct {
	GitClient *github.Client
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
		GitClient: client,
		UserContext: ctx,
	}
}

func (r *RequestHandler) ListRepositories() {
	// list all repositories for the authenticated user
	repos, _, err := r.GitClient.Repositories.List(r.UserContext, "", nil)
	if err!= nil {
		fmt.Println("Error occurred while listing repos: %s", err.Error())
	}
	for _, eachRepo := range repos {
		fmt.Println(eachRepo.Name)
		fmt.Println(eachRepo.Language)
	}
}

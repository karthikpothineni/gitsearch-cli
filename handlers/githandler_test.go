package handlers

import (
	"gopkg.in/ukautz/clif.v1"
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	"github.com/google/go-github/v38/github"
	"github.com/migueleliasweb/go-github-mock/src/mock"
	"github.com/stretchr/testify/assert"
)

var (
	accessToken = "ghp_jjkhwMWqXJzWethZdCsyuZjCAYaV072BznRD"
	orgName     = "golang"
	progressBar = clif.NewProgressBar(100).SetStyle(clif.ProgressBarStyleAscii).SetRenderWidth(80)
)

// TestNewRequestHandler - test request handler object creation
func TestNewRequestHandler(t *testing.T) {
	check := assert.New(t)

	requestHandler, err := NewRequestHandler(accessToken)
	check.Equal(nil, err)
	check.NotNil(requestHandler.GitClient)
	check.NotNil(requestHandler.UserContext)
}

// TestNewRequestHandlerEmptyToken - test request handler object creation when empty access token is passed
func TestNewRequestHandlerEmptyToken(t *testing.T) {
	check := assert.New(t)

	requestHandler, err := NewRequestHandler("")
	check.NotEqual(nil, err)
	check.Nil(requestHandler)
	check.Nil(requestHandler)
}

// TestFetchRepositoriesInfoWithOutError - test retrieving repository information when there is no error
func TestFetchRepositoriesInfoWithOutError(t *testing.T) {
	check := assert.New(t)

	// create git request handler
	requestHandler, _ := NewRequestHandler(accessToken)

	// mock git handler
	mockedHTTPClient := getMocKedHTTPClientSuccessResponse()
	mockedGitClient := github.NewClient(mockedHTTPClient)
	requestHandler.GitClient = mockedGitClient

	// fetch repo info by mocking git request handler
	contributorRepoDetails, repoLanguageDetails, err := requestHandler.FetchRepositoriesInfo(orgName, progressBar)

	check.Equal(nil, err)
	check.Equal([]string{"protobuf"}, contributorRepoDetails["gouser;;"])
	check.Equal([]string{"Go"}, repoLanguageDetails["protobuf"])
}

// TestFetchRepositoriesInfoWithError - test retrieving repository information when there is an error
func TestFetchRepositoriesInfoWithError(t *testing.T) {
	check := assert.New(t)

	// create git request handler
	requestHandler, _ := NewRequestHandler(accessToken)

	// mock git handler
	mockedHTTPClient := getMocKedHTTPClientFailureResponse()
	mockedGitClient := github.NewClient(mockedHTTPClient)
	requestHandler.GitClient = mockedGitClient

	// fetch repo info by mocking git request handler
	contributorRepoDetails, repoLanguageDetails, err := requestHandler.FetchRepositoriesInfo(orgName, progressBar)

	check.Contains(err.Error(), "error occurred while listing repos")
	check.Nil(contributorRepoDetails)
	check.Nil(repoLanguageDetails)
}

// TestHandleResponse - test if the response if printed
func TestHandleResponse(t *testing.T) {
	contributorRepoDetails := map[string][]string{
		"gouser": {"crypto"},
	}
	repoLanguageDetails := map[string][]string{
		"crypto": {"Go", "Makefile"},
	}
	rescueStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	check := assert.New(t)

	// create git request handler
	requestHandler, _ := NewRequestHandler(accessToken)

	// print response
	requestHandler.HandleResponse(contributorRepoDetails, repoLanguageDetails)

	w.Close()
	out, _ := ioutil.ReadAll(r)
	os.Stdout = rescueStdout
	check.Contains(string(out), "gouser;crypto;Go, Makefile")
}

// getMocKedHTTPClientSuccessResponse - returns the mocked http client with successful response
func getMocKedHTTPClientSuccessResponse() *http.Client {
	mockedHTTPClient := mock.NewMockedHTTPClient(
		mock.WithRequestMatch(
			mock.GetOrgsReposByOrg,
			[][]byte{
				mock.MustMarshal([]github.Repository{{
					Name: github.String("protobuf"),
				}}),
			},
		),
		mock.WithRequestMatch(
			mock.GetReposContributorsByOwnerByRepo,
			[][]byte{
				mock.MustMarshal([]github.Contributor{{
					Login: github.String("gouser"),
				}}),
			},
		),
		mock.WithRequestMatch(
			mock.GetReposLanguagesByOwnerByRepo,
			[][]byte{
				mock.MustMarshal(map[string]int{
					"Go": 486640,
				}),
			},
		),
	)

	return mockedHTTPClient
}

// getMocKedHTTPClientFailureResponse - returns the mocked http client with failure response
func getMocKedHTTPClientFailureResponse() *http.Client {
	mockedHTTPClient := mock.NewMockedHTTPClient(
		mock.WithRequestMatchHandler(
			mock.GetOrgsReposByOrg,
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				mock.WriteError(
					w,
					http.StatusNotFound,
					"Invalid organization name",
				)
			}),
		),
	)

	return mockedHTTPClient
}

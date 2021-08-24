# gitsearch-cli
CLI for fetching github repository information

## Description
This application is responsible for retrieving github repository information which includes contributors, languages for a repository. Internally, this application uses "https://github.com/ukautz/clif" and "https://github.com/google/go-github/" libraries.

## Setup
### Local
1. Clone the repository under GOPATH
2. Install dependencies using ```go mod download```
3. Run application using 

```go run report.go report --organization {org name} --auth-key {personal token}```

**Note:**
Replace ```{org name}``` with github organization name and ```{personal token}``` with github [personal access token](https://docs.github.com/en/github/authenticating-to-github/keeping-your-account-and-data-secure/creating-a-personal-access-token)
### Run Linter
```golangci-lint run -v -c golangci.yml```
### Run Tests
```go test -v -cover ./...```

## Example

![](https://github.com/karthikpothineni/staticfiles/blob/main/gif/gitsearch-cli.gif)

### Sample Output
```$xslt
login;name;email;repositories;languages
avelino;;;lint, gofrontend;Makefile, Awk, C++, C, M4, Assembly, Fortran, Go, Shell
quasilyte;;;protobuf;Go, Shell
spf13;;;blog;
dddent;;;image;Go, Assembly
tombergan;;;net;Makefile, Go, Dockerfile, Assembly
pasztorpisti;;;mock;Go, Shell
crawshaw;;;lint, go, review, exp, mobile, talks, build;CSS, Shell, Go, Makefile, HTML, Dockerfile
```

## Code Coverage
Current code coverage is more than **85%**

## Circle CI
Added Circle CI integration. For every commit, both the lint and tests will be executed
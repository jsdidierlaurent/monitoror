//go:generate mockery --name SearchService

package gogithub

import (
	"context"

	githubApi "github.com/google/go-github/github"
)

type SearchService interface {
	Issues(ctx context.Context, query string, opt *githubApi.SearchOptions) (*githubApi.IssuesSearchResult, *githubApi.Response, error)
}

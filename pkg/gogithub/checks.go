//go:generate mockery --name ChecksService

package gogithub

import (
	"context"

	githubApi "github.com/google/go-github/github"
)

type ChecksService interface {
	ListCheckRunsForRef(ctx context.Context, owner, repo, ref string, opt *githubApi.ListCheckRunsOptions) (*githubApi.ListCheckRunsResults, *githubApi.Response, error)
}

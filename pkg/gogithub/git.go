//go:generate mockery --name GitService

package gogithub

import (
	"context"

	githubApi "github.com/google/go-github/github"
)

type GitService interface {
	GetCommit(ctx context.Context, owner string, repo string, sha string) (*githubApi.Commit, *githubApi.Response, error)
}

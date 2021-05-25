//go:generate mockery --name RepositoriesService

package gogithub

import (
	"context"

	githubApi "github.com/google/go-github/github"
)

type RepositoriesService interface {
	ListStatuses(ctx context.Context, owner, repo, ref string, opt *githubApi.ListOptions) ([]*githubApi.RepoStatus, *githubApi.Response, error)
}

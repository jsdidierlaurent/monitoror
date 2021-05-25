//go:generate mockery --name Repository

package api

import "github.com/monitoror/monitoror/monitorables/github/api/models"

type (
	Repository interface {
		GetCount(query string) (int, error)
		GetChecks(owner, repository, ref string) (*models.Checks, error)
		GetPullRequest(owner, repository string, id int) (*models.PullRequest, error)
		GetPullRequests(owner, repository string) ([]models.PullRequest, error)
		GetCommit(owner, repository, sha string) (*models.Commit, error)
	}
)

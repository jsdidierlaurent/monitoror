//go:generate mockery --name Repository

package api

import (
	"github.com/monitoror/monitoror/monitorables/travisci/api/models"
)

type (
	Repository interface {
		GetLastBuildStatus(owner, repository, branch string) (*models.Build, error)
	}
)

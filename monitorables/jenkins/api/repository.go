//go:generate mockery --name Repository

package api

import (
	"github.com/monitoror/monitoror/monitorables/jenkins/api/models"
)

type (
	Repository interface {
		GetJob(jobName string, branch string) (*models.Job, error)
		GetLastBuildStatus(job *models.Job) (*models.Build, error)
	}
)

//go:generate mockery --name Jenkins

package gojenkins

import (
	gojenkins "github.com/jsdidierlaurent/golang-jenkins"
)

type Jenkins interface {
	GetJob(jobName string) (job gojenkins.Job, err error)
	GetBuildByJobId(jobID string, number int) (build gojenkins.Build, err error)
	GetLastBuildByJobId(jobID string) (build gojenkins.Build, err error)
}

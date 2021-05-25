//go:generate mockery --name PipelinesService

package gogitlab

import (
	"github.com/xanzy/go-gitlab"
)

type PipelinesService interface {
	GetPipeline(pid interface{}, pipeline int, options ...gitlab.RequestOptionFunc) (*gitlab.Pipeline, *gitlab.Response, error)
	ListProjectPipelines(pid interface{}, opt *gitlab.ListProjectPipelinesOptions, options ...gitlab.RequestOptionFunc) ([]*gitlab.PipelineInfo, *gitlab.Response, error)
}

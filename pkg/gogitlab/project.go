//go:generate mockery --name ProjectService

package gogitlab

import (
	"github.com/xanzy/go-gitlab"
)

type ProjectService interface {
	GetProject(pid interface{}, opt *gitlab.GetProjectOptions, options ...gitlab.RequestOptionFunc) (*gitlab.Project, *gitlab.Response, error)
}

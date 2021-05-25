//go:generate mockery --name IssuesService

package gogitlab

import (
	"github.com/xanzy/go-gitlab"
)

type IssuesService interface {
	ListIssues(opt *gitlab.ListIssuesOptions, options ...gitlab.RequestOptionFunc) ([]*gitlab.Issue, *gitlab.Response, error)
	ListProjectIssues(pid interface{}, opt *gitlab.ListProjectIssuesOptions, options ...gitlab.RequestOptionFunc) ([]*gitlab.Issue, *gitlab.Response, error)
}

//go:generate mockery --name TravisCI

package gotravis

import (
	"context"
	"net/http"

	"github.com/shuheiktgw/go-travis"
)

type TravisCI interface {
	ListByRepoSlug(ctx context.Context, repoSlug string, opt *travis.BuildsByRepoOption) ([]*travis.Build, *http.Response, error)
}

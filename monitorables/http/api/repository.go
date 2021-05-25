//go:generate mockery --name Repository

package api

import (
	"github.com/monitoror/monitoror/monitorables/http/api/models"
)

type (
	Repository interface {
		Get(url string) (*models.Response, error)
	}
)

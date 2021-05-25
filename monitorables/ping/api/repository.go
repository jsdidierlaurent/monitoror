//go:generate mockery --name Repository

package api

import (
	"github.com/monitoror/monitoror/monitorables/ping/api/models"
)

type (
	Repository interface {
		ExecutePing(hostname string) (*models.Ping, error)
	}
)

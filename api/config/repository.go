//go:generate mockery --name Repository

package config

import (
	"github.com/monitoror/monitoror/api/config/models"
)

type (
	Repository interface {
		GetConfigFromURL(url string) (*models.Config, error)
		GetConfigFromPath(baseDir, filePath string) (*models.Config, error)
	}
)

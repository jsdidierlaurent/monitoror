//go:generate mockery --name Usecase

package api

import (
	coreModels "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorables/port/api/models"
)

const (
	PortTileType coreModels.TileType = "PORT"
)

type (
	Usecase interface {
		Port(params *models.PortParams) (*coreModels.Tile, error)
	}
)

//go:generate mockery --name Usecase

package api

import (
	coreModels "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorables/ping/api/models"
)

const (
	PingTileType coreModels.TileType = "PING"
)

type (
	Usecase interface {
		Ping(params *models.PingParams) (*coreModels.Tile, error)
	}
)

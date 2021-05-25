//go:generate mockery --name Usecase

package api

import (
	coreModels "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorables/travisci/api/models"
)

const (
	TravisCIBuildTileType coreModels.TileType = "TRAVISCI-BUILD"
)

type (
	Usecase interface {
		Build(params *models.BuildParams) (*coreModels.Tile, error)
	}
)

//go:generate mockery --name Usecase

package api

import (
	coreModels "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorables/azuredevops/api/models"
)

const (
	AzureDevOpsBuildTileType   coreModels.TileType = "AZUREDEVOPS-BUILD"
	AzureDevOpsReleaseTileType coreModels.TileType = "AZUREDEVOPS-RELEASE"
)

type (
	Usecase interface {
		Build(params *models.BuildParams) (*coreModels.Tile, error)
		Release(params *models.ReleaseParams) (*coreModels.Tile, error)
	}
)

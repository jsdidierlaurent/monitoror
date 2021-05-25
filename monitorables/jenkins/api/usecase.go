//go:generate mockery --name Usecase

package api

import (
	uiConfigModels "github.com/monitoror/monitoror/api/config/models"
	coreModels "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorables/jenkins/api/models"
)

const (
	JenkinsBuildTileType coreModels.TileType = "JENKINS-BUILD"
)

type (
	Usecase interface {
		Build(params *models.BuildParams) (*coreModels.Tile, error)
		BuildGenerator(params interface{}) ([]uiConfigModels.GeneratedTile, error)
	}
)

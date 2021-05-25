//go:generate mockery --name Usecase

package api

import (
	uiConfigModels "github.com/monitoror/monitoror/api/config/models"
	coreModels "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorables/github/api/models"
)

const (
	GithubCountTileType       coreModels.TileType = "GITHUB-COUNT"
	GithubChecksTileType      coreModels.TileType = "GITHUB-CHECKS"
	GithubPullRequestTileType coreModels.TileType = "GITHUB-PULLREQUEST"
)

type (
	Usecase interface {
		Count(params *models.CountParams) (*coreModels.Tile, error)
		Checks(params *models.ChecksParams) (*coreModels.Tile, error)
		PullRequest(params *models.PullRequestParams) (*coreModels.Tile, error)

		PullRequestsGenerator(params interface{}) ([]uiConfigModels.GeneratedTile, error)
	}
)

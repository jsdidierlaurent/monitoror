package usecase

import (
	"fmt"

	"github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorable/stripe"
	stripeModels "github.com/monitoror/monitoror/monitorable/stripe/models"
	"github.com/monitoror/monitoror/pkg/monitoror/cache"
)

const buildCacheSize = 5

type (
	stripeUsecase struct {
		repository stripe.Repository

		// builds cache. used for save small history of build for stats
		buildsCache *cache.BuildCache
	}
)

func NewStripeUsecase(repository stripe.Repository) stripe.Usecase {
	return &stripeUsecase{
		repository,
		cache.NewBuildCache(buildCacheSize),
	}
}

func (su *stripeUsecase) Count(params *stripeModels.CountParams) (*models.Tile, error) {
	tile := models.NewTile(stripe.StripeCountTileType).WithValue(models.RawUnit)
	tile.Label = "Stripe Count"

	net, count, err := su.repository.GetCount(params.CreatedAfter)
	if err != nil {
		return nil, &models.MonitororError{Err: err, Tile: tile, Message: "timestamp has invalid format"}
	}
	tile.Status = models.SuccessStatus
	tile.Value.Values = append(tile.Value.Values, fmt.Sprintf("$%.2f (%d)", net, count))
	return tile, nil
}

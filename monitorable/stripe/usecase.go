package stripe

import (
	"github.com/monitoror/monitoror/models"
	stripeModels "github.com/monitoror/monitoror/monitorable/stripe/models"
)

const (
	StripeCountTileType models.TileType = "STRIPE-COUNT"
)

type (
	Usecase interface {
		Count(params *stripeModels.CountParams) (*models.Tile, error)
	}
)

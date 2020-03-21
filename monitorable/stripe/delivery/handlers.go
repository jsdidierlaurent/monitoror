package delivery

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorable/stripe"
	stripeModels "github.com/monitoror/monitoror/monitorable/stripe/models"
)

type StripeDelivery struct {
	stripeUsecase stripe.Usecase
}

func NewStripeDelivery(u stripe.Usecase) *StripeDelivery {
	return &StripeDelivery{u}
}

func (d *StripeDelivery) GetCount(c echo.Context) error {
	// Bind / check Params
	params := &stripeModels.CountParams{}
	err := c.Bind(params)
	if err != nil || !params.IsValid() {
		return models.QueryParamsError
	}

	tile, err := d.stripeUsecase.Count(params)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, tile)
}

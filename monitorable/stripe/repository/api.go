package repository

import (
	"strconv"
	"time"

	"github.com/monitoror/monitoror/config"
	stripeLocal "github.com/monitoror/monitoror/monitorable/stripe"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/client"
)

type (
	stripeRepository struct {
		stripeClient *client.API
		config       *config.Stripe
	}
)

// NewStripeRepository makes a new Stripe connection from an API key
func NewStripeRepository(config *config.Stripe) stripeLocal.Repository {
	stripe.DefaultLeveledLogger = &stripe.LeveledLogger{Level: stripe.LevelError}
	sc := &client.API{}
	sc.Init(config.Token, nil)
	return &stripeRepository{
		stripeClient: sc,
		config:       config,
	}
}

func (r *stripeRepository) GetCount(afterTimestamp string) (float64, int, error) {
	if afterTimestamp == "" || afterTimestamp == "today" {
		afterTimestamp = strconv.FormatInt(bod(time.Now().Local()).Unix(), 10)
	}
	_, err := strconv.ParseUint(afterTimestamp, 10, 64)
	if err != nil {
		return 0, 0, err
	}
	params := &stripe.BalanceTransactionListParams{}
	params.Filters.AddFilter("type", "", "charge")
	params.Filters.AddFilter("created", "gte", afterTimestamp)
	params.Filters.AddFilter("limit", "", "100")
	result := r.stripeClient.BalanceTransaction.List(params)
	count := 0
	net := 0.0
	for result.Next() {
		net = net + float64(result.BalanceTransaction().Net)
		count = count + 1
	}
	net = net / 100
	return net, count, nil
}

func bod(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, t.Location())
}

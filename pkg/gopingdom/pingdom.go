//go:generate mockery --name PingdomCheckAPI|PingdomTransactionCheckAPI

package gopingdom

import "github.com/jsdidierlaurent/go-pingdom/pingdom"

type PingdomCheckAPI interface {
	List(params ...map[string]string) ([]pingdom.CheckResponse, error)
	Read(id int) (*pingdom.CheckResponse, error)
}

type PingdomTransactionCheckAPI interface {
	List(params ...map[string]string) ([]pingdom.TransactionCheckResponse, error)
	Read(id int) (*pingdom.TransactionCheckResponse, error)
}

//go:generate mockery --name Repository

package api

type (
	Repository interface {
		OpenSocket(hostname string, port int) error
	}
)

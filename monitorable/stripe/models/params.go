package models

type (
	CountParams struct {
		CreatedAfter string
	}
)

func (p *CountParams) IsValid() bool {
	return true
}

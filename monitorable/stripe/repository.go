package stripe

type (
	Repository interface {
		GetCount(afterTimestamp string) (float64, int, error)
	}
)

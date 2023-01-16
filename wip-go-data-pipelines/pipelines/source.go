package pipelines

import "context"

type Pagination struct {
	PageNumber int
	PageSize   int
}

// ListNumbers emulates a typical data source creating a fictitious result considering the Pagination parameter
func ListNumbers(_ context.Context, p Pagination) ([]int, error) {
	result := make([]int, 0, p.PageSize)

	start := (p.PageNumber - 1) * p.PageSize
	finish := p.PageNumber * p.PageSize

	for i := start; i < finish; i++ {
		result = append(result, i)
	}
	return result, nil
}

func NumbersCount(_ context.Context) (int, error) {
	return 100, nil
}

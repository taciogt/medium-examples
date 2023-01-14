package pipelines

type Pagination struct {
	Start int
	Limit int
}

func ListNumbers(p Pagination) ([]int, error) {
	result := make([]int, 0, p.Limit)

	for i := p.Start; i < p.Start+p.Limit; i++ {
		result = append(result, i)
	}
	return result, nil
}

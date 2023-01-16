package pipelines

import "context"

func Pipeline(ctx context.Context) (chan int, chan error) {
	resultInt := make(chan int)
	resultErr := make(chan error)

	go func() {
		defer close(resultInt)
		defer close(resultErr)

		totalNumbers, err := NumbersCount(ctx)
		if err != nil {
			resultErr <- err
			return
		}

		pageSize := 10
		pagesCount := totalNumbers / pageSize

		for pageNumber := 1; pageNumber <= pagesCount; pageNumber++ {
			ns, err := ListNumbers(ctx, Pagination{
				PageNumber: pageNumber,
				PageSize:   pageSize,
			})
			if err != nil {
				resultErr <- err
				return
			}
			for _, n := range ns {
				select {
				case resultInt <- n:
				case <-ctx.Done():
					resultErr <- ctx.Err()
					return
				}
			}
		}
	}()

	return resultInt, resultErr
}

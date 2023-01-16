package pipelines

import (
	"context"
	"golang.org/x/sync/errgroup"
)

func Pipeline(ctx context.Context) (chan int, chan error) {
	resultInt := make(chan int)
	resultErr := make(chan error)

	go func() {
		defer close(resultInt)
		defer close(resultErr)

		g, ctx := errgroup.WithContext(ctx)

		totalNumbers, err := NumbersCount(ctx)
		if err != nil {
			resultErr <- err
			return
		}

		pageSize := 10
		pagesCount := totalNumbers / pageSize

		for pageNumber := 1; pageNumber <= pagesCount; pageNumber++ {
			g.Go(func() error {
				ns, err := ListNumbers(ctx, Pagination{
					PageNumber: pageNumber,
					PageSize:   pageSize,
				})
				if err != nil {
					return err
				}
				for _, n := range ns {
					select {
					case resultInt <- n:
					case <-ctx.Done():
						return ctx.Err()
					}
				}
				return nil
			})
		}

		if err := g.Wait(); err != nil {
			resultErr <- err
		}
	}()

	return resultInt, resultErr
}

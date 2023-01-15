package pipelines

import "context"

func Pipeline(ctx context.Context) (chan int, chan error) {
	resultInt := make(chan int)
	resultErr := make(chan error)

	go func() {
		ns, err := ListNumbers(ctx, Pagination{
			PageNumber: 1,
			PageSize:   10,
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
	}()

	return resultInt, resultErr
}

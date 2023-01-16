package pipelines

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestPipeline(t *testing.T) {
	want := make([]int, 100)
	for i := 0; i < 100; i++ {
		want[i] = i
	}

	tests := []struct {
		name    string
		want    []int
		wantErr error
	}{{
		name: "first test",
		want: want,
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			gotIntCh, gotErrCh := Pipeline(ctx)
			gotInt := make([]int, 0)
			var gotErr error

			go func() {
				for n := range gotIntCh {
					gotInt = append(gotInt, n)
				}
				cancel()
			}()

			select {
			case gotErr = <-gotErrCh:
			case <-ctx.Done():
			}

			if err := ctx.Err(); !errors.Is(err, context.Canceled) {
				t.Errorf("Pipeline() didn't close results channel before timeout reached. err=%v", err)
				return
			}

			if len(gotInt) != len(tt.want) {
				t.Errorf("Pipeline() returned results quantity different than expected. got=%v, want=%v",
					len(gotInt), len(tt.want))
			}

			if gotErr != tt.wantErr {
				t.Errorf("Pipeline() gotErr = %v, want %v", gotErr, tt.wantErr)
			}
		})
	}
}

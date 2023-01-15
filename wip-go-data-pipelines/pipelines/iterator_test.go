package pipelines

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"
)

func TestPipeline(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name  string
		want  []int
		want1 chan error
	}{
		{
			name:  "first test",
			want:  []int{0, 1, 2, 3, 4, 5, 6, 7},
			want1: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			gotIntCh, gotErrCh := Pipeline(ctx)
			gotInt := make([]int, 0)

			go func() {
				for n := range gotIntCh {
					gotInt = append(gotInt, n)
				}
				cancel()
			}()
			select {
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
			if !reflect.DeepEqual(gotErrCh, tt.want1) {
				t.Errorf("Pipeline() gotErrCh = %v, want %v", gotErrCh, tt.want1)
			}
		})
	}
}

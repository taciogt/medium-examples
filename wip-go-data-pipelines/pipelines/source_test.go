package pipelines

import (
	"context"
	"reflect"
	"testing"
)

func TestListNumbers(t *testing.T) {
	type args struct {
		p Pagination
	}
	tests := []struct {
		name       string
		pagination Pagination
		want       []int
		wantErr    bool
	}{{
		name: "simplest case",
		pagination: Pagination{
			PageNumber: 1,
			PageSize:   10,
		},
		want: []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
	},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ListNumbers(context.TODO(), tt.pagination)
			if (err != nil) != tt.wantErr {
				t.Errorf("ListNumbers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ListNumbers() got = %v, want %v", got, tt.want)
			}
		})
	}
}

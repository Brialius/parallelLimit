package parallelLimit

import (
	"errors"
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestParallelLimit(t *testing.T) {
	type args struct {
		funcs      []func() error
		maxWorkers int
		maxErrors  int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"No errors/maxWorkers 2/maxErrors 0",
			args{
				[]func() error{
					taskFunc(1, nil),
					taskFunc(2, nil),
					taskFunc(3, nil),
					taskFunc(4, nil),
					taskFunc(5, nil),
					taskFunc(6, nil),
					taskFunc(7, nil),
					taskFunc(8, nil),
					taskFunc(9, nil),
					taskFunc(10, nil),
				},
				2,
				0,
			},
			false,
		},
		{
			"Acceptable quantity of errors/maxWorkers 4/maxErrors 2",
			args{
				[]func() error{
					taskFunc(1, nil),
					taskFunc(2, nil),
					taskFunc(3, errors.New("error")),
					taskFunc(4, nil),
					taskFunc(5, nil),
					taskFunc(6, errors.New("error")),
					taskFunc(7, nil),
					taskFunc(8, nil),
					taskFunc(9, nil),
					taskFunc(10, nil),
				},
				4,
				2,
			},
			false,
		},
		{
			"Too much errors/maxWorkers 5/maxErrors 3",
			args{
				[]func() error{
					taskFunc(1, nil),
					taskFunc(2, errors.New("error")),
					taskFunc(3, errors.New("error")),
					taskFunc(4, errors.New("error")),
					taskFunc(5, errors.New("error")),
					taskFunc(6, nil),
					taskFunc(7, nil),
					taskFunc(8, errors.New("error")),
					taskFunc(9, nil),
					taskFunc(10, nil),
				},
				5,
				3,
			},
			true,
		},
	}
	rand.Seed(time.Now().UnixNano())
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ParallelLimit(tt.args.funcs, tt.args.maxWorkers, tt.args.maxErrors); (err != nil) != tt.wantErr {
				t.Errorf("ParallelLimit() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func taskFunc(i int, err error) func() error {
	return func() error {
		if err != nil {
			fmt.Printf("task %d failed\n", i)
			return err
		}
		time.Sleep(time.Second * time.Duration(rand.Intn(3)+1))
		fmt.Printf("task %d completed\n", i)
		return nil
	}
}

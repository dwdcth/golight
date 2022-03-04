package task

import (
	"fmt"
	"math"
	"math/rand"
	"testing"
	"time"
)

func TestWatchTickTask(t *testing.T) {
	type args struct {
		watchInterval   time.Duration
		getIntervalFunc GetIntervalFunc
		taskFunc        TaskFunc
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test",
			args: args{
				watchInterval: 5 * time.Second,
				getIntervalFunc: func() (time.Duration, error) {
					interval := []int{2, 4, 3, 3, 2, 2, 4}
					i := math.Floor(rand.Float64() * float64(len(interval)))
					if i == 0 {
						i = 2
					}
					return time.Second * time.Duration(i), nil
				},
				taskFunc: func() error {
					fmt.Println("********stats-start")
					fmt.Println("现在开始等待3秒,time=", time.Now().Unix())
					time.Sleep(1 * time.Second)
					fmt.Println("********stats-over")
					return nil
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			WatchTickTask(tt.args.watchInterval, tt.args.getIntervalFunc, tt.args.taskFunc)
		})
	}
}

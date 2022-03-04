package task

import (
	"context"
	"log"
	"time"
)

type a = int

type GetIntervalFunc func() (time.Duration, error)
type TaskFunc func() error

// WatchTickTask 定时监控任务执行时间变化，并按照变化后的时间定时执行任务
// watchInterval: 监控间隔
// getIntervalFunc: 获取间隔时间的函数
// taskFunc: 任务函数
// 注意如果监控的时间间隔比任务运行的时间小，任务的时间每次都变化的情况下
// 并且统计任务运行的时间比监控的时间长，那么统计任务会出现经常被取消的情况
func WatchTickTask(watchInterval time.Duration, getIntervalFunc GetIntervalFunc, taskFunc TaskFunc) {
	if watchInterval <= 0 {
		log.Panic("watchInterval must > 0")
	}
	queryTicker := time.NewTicker(watchInterval)
	defer queryTicker.Stop()

	var (
		ctx      context.Context
		cancel   context.CancelFunc
		interval time.Duration
		err      error
	)
	ctx, cancel = context.WithCancel(context.Background())
	for t := range queryTicker.C {
		log.Printf("************************query ticker %s************************\n", t)
		// 第一次运行
		if interval == 0 {
			interval, err = getIntervalFunc()
			if err != nil {
				log.Printf("get interval error: %s\n", err)
				continue
			}
			if interval <= 0 {
				log.Println("interval must > 0")
				continue
			}
			go runTask(ctx, interval, taskFunc)
			continue
		}

		newInterval, err := getIntervalFunc()
		if err != nil {
			log.Printf("get interval error: %s\n", err)
			continue
		}
		if newInterval <= 0 {
			log.Println("interval must > 0")
			continue
		}
		//若新查询时间间隔与历史的不一致，则结束父级协程
		if newInterval != interval {
			interval = newInterval
			cancel() //发送结束信号
			ctx, cancel = context.WithCancel(context.Background())
			go runTask(ctx, interval, taskFunc)
		}

	}
	cancel()
}

func runTask(ctx context.Context, interval time.Duration, taskFunc TaskFunc) {
	statsTicker := time.NewTicker(interval)
	defer statsTicker.Stop()
	for {
		select {
		case <-statsTicker.C:
			err := taskFunc()
			if err != nil {
				log.Printf("task error: %s\n", err)
			}
		case <-ctx.Done():
			log.Println("task interval done")
			return
		}
	}
}

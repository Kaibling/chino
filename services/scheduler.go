package services

import (
	"context"
	"fmt"
	"time"

	"chino/pkg/log"
)

type SchedulerObject interface {
	Run() error
}

type SchedulerService struct {
	o   SchedulerObject
	ctx context.Context
}

func NewSchedulerService(ctx context.Context, o SchedulerObject) *SchedulerService {
	return &SchedulerService{o: o, ctx: ctx}
}

func (s *SchedulerService) Start(t time.Duration) chan bool {
	done := make(chan bool)
	Runf(s.ctx, t, done, s.o.Run)
	return done
}

func Runf(ctx context.Context, t time.Duration, done chan bool, f func() error) {
	ticker := time.NewTicker(t)
	go func() {
		for {
			select {
			case <-done:
				ticker.Stop()
				log.Info(ctx, "scheduler shutdown")
				return
			case <-ticker.C:
				log.Info(ctx, fmt.Sprintf("run at %s", time.Now()))
				err := f()
				if err != nil {
					log.Error(ctx, err)
					log.Info(ctx, "Scheduled Function stopped")
					return
				}
				log.Info(ctx, fmt.Sprintf("run finished at %s", time.Now()))
			}
		}
	}()
}

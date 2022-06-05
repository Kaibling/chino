package services

import (
	"chino/lib/scheduler"
	"time"
)

type SchedulerObject interface {
	Run() error
}

type SchedulerService struct {
	o SchedulerObject
}

func NewSchedulerService(o SchedulerObject) *SchedulerService {
	return &SchedulerService{o}
}

func (s *SchedulerService) Startf(t time.Duration, f func() error) chan bool {
	done := make(chan bool)
	scheduler.Runf(t, done, f)
	return done
}

func (s *SchedulerService) Start(t time.Duration) chan bool {
	done := make(chan bool)
	scheduler.Runf(t, done, s.o.Run)
	return done
}

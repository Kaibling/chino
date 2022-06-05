package scheduler

import (
	"chino/lib/logging"
	"fmt"
	"time"
)

func Runf(t time.Duration, done chan bool, f func() error) {
	ticker := time.NewTicker(t)
	go func() {
		for {
			select {
			case <-done:
				ticker.Stop()
				fmt.Println("scheduler shutdown")
				return
			case <-ticker.C:
				logging.Logger.Infof("run at %s", time.Now())
				err := f()
				if err != nil {
					logging.Logger.Error(err.Error())
					logging.Logger.Infof("Scheduled Function stopped")
					return
				}
				logging.Logger.Infof("run finished at %s", time.Now())
			}
		}
	}()
}

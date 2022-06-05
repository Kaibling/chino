package main

import (
	"chino/api"
	"chino/gormrepo"
	"chino/gormrepo/database"
	"chino/notification"
	"chino/services"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	db, err := database.InitDatabase()
	if err != nil {
		fmt.Println(err.Error())
	}
	err = gormrepo.Migrate(db)
	if err != nil {
		fmt.Println(err.Error())
	}
	movieRepo := gormrepo.NewMovieRepo(db)
	movieService := services.NewMovieService(movieRepo)
	notifier := notification.NewFmtNotifier()
	notificationService := services.NewNotificationService(notifier)
	movieService.AddNotificationService(notificationService)
	schedulerService := services.NewSchedulerService(movieService)
	cancel := schedulerService.Start(5 * time.Second)

	c := make(chan os.Signal, 1)
	apiServerDone := make(chan bool, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go api.NewServer(db, apiServerDone)

	<-c
	apiServerDone <- true
	cancel <- true
	fmt.Println("closing time")
	time.Sleep(2 * time.Second)
	os.Exit(1)

}

package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"chino/api"
	"chino/models"
	"chino/pkg/config"
	"chino/pkg/log"
	"chino/pkg/notification"
	"chino/pkg/persistence/repo_sqlx"
	"chino/services"

	goflag "flag"

	"github.com/jmoiron/sqlx"
	flag "github.com/spf13/pflag"
)

var appVersion string
var goVersion string
var buildTime string

func main() {
	showVersion := flag.BoolP("version", "v", false, "help app version")
	flag.CommandLine.AddGoFlagSet(goflag.CommandLine)
	flag.Parse()
	if *showVersion {
		fmt.Printf("app version: %s\n", appVersion)
		fmt.Printf("go version: %s\n", goVersion)
		fmt.Printf("build date: %s\n", buildTime)
		return
	}

	cfg := config.Init()

	l, err := log.New()
	if err != nil {
		fmt.Println(err.Error())
	}
	ctx := context.WithValue(context.Background(), models.String("logger"), l)

	db, err := repo_sqlx.InitDatabase(*cfg)
	if err != nil {
		log.Error(ctx, err)
		return
	}

	err = repo_sqlx.InitDB(ctx, db)
	if err != nil {
		log.Error(ctx, err)
	}

	schedulerService := InitScheduler(ctx, db)
	cancel := schedulerService.Start(5 * time.Hour * 24)

	c := make(chan os.Signal, 1)
	apiServerDone := make(chan bool, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go api.NewServer(ctx, db, *cfg, l, apiServerDone)

	<-c
	apiServerDone <- true
	cancel <- true
	log.Info(ctx, "closing time")
	time.Sleep(2 * time.Second)
	os.Exit(1)
}

func InitScheduler(ctx context.Context, db *sqlx.DB) *services.SchedulerService {
	movieRepo := repo_sqlx.NewMovieRepo(ctx, db)
	movieService := services.NewMovieService(ctx, movieRepo)
	notifier := notification.NewFmtNotifier(ctx)
	notificationService := services.NewNotificationService(notifier)
	movieService.AddNotificationService(notificationService)
	return services.NewSchedulerService(ctx, movieService)
}

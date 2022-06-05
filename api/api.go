package api

import (
	"chino/api/dashboard"
	"chino/api/movies"
	"chino/crawler/uncut"
	"chino/gormrepo"
	"chino/services"
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"gorm.io/gorm"
)

func NewServer(db *gorm.DB, done chan bool) {

	movieService := services.NewMovieService(gormrepo.NewMovieRepo(db))
	crawlerService := services.NewCrawlerService(uncut.NewCrawler())

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(injectData("movieservice", movieService))
	r.Use(injectData("crawlerservice", crawlerService))
	r.Route("/", func(r chi.Router) {
		r.Mount("/movies", movies.BuildRoute())
		r.Mount("/dashboard", dashboard.BuildRoute())
	})
	server := http.Server{Addr: ":3333", Handler: r}
	serverCtx, serverStopCtx := context.WithCancel(context.Background())
	go func() {
		<-done
		fmt.Println("shutown api server")
		// Shutdown signal with grace period of 30 seconds
		shutdownCtx, _ := context.WithTimeout(serverCtx, 30*time.Second)

		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				log.Fatal("graceful shutdown timed out.. forcing exit.")
			}
		}()

		// Trigger graceful shutdown
		err := server.Shutdown(shutdownCtx)
		if err != nil {
			log.Fatal(err)
		}
		serverStopCtx()
	}()
	fmt.Println("listening on 3000")
	err := http.ListenAndServe(":3000", r)
	if err != nil {
		fmt.Println(err.Error())
	}

	return
}

func injectData(key string, data interface{}) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), key, data)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}
}

package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/Kaibling/chino/api/dashboard"
	"github.com/Kaibling/chino/api/health"
	"github.com/Kaibling/chino/api/movies"
	"github.com/Kaibling/chino/models"
	"github.com/Kaibling/chino/pkg/api_middleware"
	"github.com/Kaibling/chino/pkg/config"
	"github.com/Kaibling/chino/pkg/log"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jmoiron/sqlx"
	"moul.io/chizap"
)

func NewServer(ctx context.Context, db *sqlx.DB, c config.Config, l *log.Logger, done chan bool) {
	listeningStr := c.BindingIP + ":" + c.BindingPort
	r := chi.NewRouter()

	r.Use(injectData("logger", l))
	r.Use(chizap.New(l.Raw(), &chizap.Opts{
		WithReferer:   true,
		WithUserAgent: true,
	}))
	r.Use(middleware.Recoverer)

	r.Use(injectData("db", db))
	r.Use(api_middleware.AddEnvelope)
	r.Route("/", func(r chi.Router) {
		r.Mount("/movies", movies.BuildRoute())
		r.Mount("/dashboard", dashboard.BuildRoute())
		r.Mount("/health", health.BuildRoute())
	})
	server := http.Server{Addr: listeningStr, Handler: r}
	serverCtx, serverStopCtx := context.WithCancel(context.Background())
	go func() {
		<-done
		log.Info(ctx, "shutown api server")
		// Shutdown signal with grace period of 30 seconds
		shutdownCtx, cancel := context.WithTimeout(serverCtx, 5*time.Second)

		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				log.Warn(ctx, "graceful shutdown timed out.. forcing exit.")
			}
		}()

		// Trigger graceful shutdown
		err := server.Shutdown(shutdownCtx)
		if err != nil {
			log.Error(ctx, err)
		}
		serverStopCtx()
		cancel()
	}()
	log.Info(ctx, fmt.Sprintf("listening on %s", listeningStr))
	err := http.ListenAndServe(listeningStr, r)
	if err != nil {
		log.Error(ctx, err)
	}
}

func injectData(key string, data interface{}) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), models.String(key), data)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}
}

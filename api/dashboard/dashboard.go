package dashboard

import (
	"net/http"
	"strconv"
	"text/template"
	"time"

	"github.com/Kaibling/chino/crawler/uncut"
	"github.com/Kaibling/chino/pkg/log"
	"github.com/Kaibling/chino/pkg/persistence/repo_sqlx"
	"github.com/Kaibling/chino/pkg/utils"
	"github.com/Kaibling/chino/services"

	"github.com/jmoiron/sqlx"
)

func extract(r *http.Request) (*services.MovieService, *services.CrawlerService) {
	db := utils.GetContext("db", r).(*sqlx.DB)
	mr := repo_sqlx.NewMovieRepo(r.Context(), db)
	ms := services.NewMovieService(r.Context(), mr)

	cr := uncut.NewCrawler(r.Context())
	cs := services.NewCrawlerService(cr)
	return ms, cs
}

func show(w http.ResponseWriter, r *http.Request) {
	ms, cs := extract(r)
	var month int
	key, ok := r.URL.Query()["month"]
	if !ok {
		_, m, _ := time.Now().Date()
		month = int(m)

	} else {
		m, err := strconv.Atoi(key[0])
		if err != nil {
			log.Error(r.Context(), err)
			return
		}
		month = m
	}

	objects, err := cs.GetMovies(month)
	if err != nil {
		log.Error(r.Context(), err)
		return
	}
	favorites, err := ms.ReadUntil(month)
	if err != nil {
		log.Error(r.Context(), err)
		return
	}
	movieList := ms.FavoriteMovies(favorites, objects)
	tmpl := template.Must(template.ParseFiles("api/dashboard/dashboard.html"))
	if err := tmpl.Execute(w, movieList); err != nil {
		log.Error(r.Context(), err)
	}
}

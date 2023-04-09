package movies

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/Kaibling/chino/models"
	"github.com/Kaibling/chino/pkg/log"
	"github.com/Kaibling/chino/pkg/persistence/repo_sqlx"
	"github.com/Kaibling/chino/pkg/utils"
	"github.com/Kaibling/chino/services"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
)

func extract(r *http.Request) (*utils.Envelope, *services.MovieService) {
	envelope := utils.GetContext("envelope", r).(*utils.Envelope)

	db := utils.GetContext("db", r).(*sqlx.DB)
	mr := repo_sqlx.NewMovieRepo(r.Context(), db)
	ms := services.NewMovieService(r.Context(), mr)
	return envelope, ms
}

func create(w http.ResponseWriter, r *http.Request) {
	envelope, rs := extract(r)
	var model models.Movie
	err := json.NewDecoder(r.Body).Decode(&model)
	if err != nil {
		log.Error(r.Context(), err)
		utils.Render(w, r, envelope.SetError(err))
		return
	}
	m, err := rs.Create(model)
	if err != nil {
		log.Error(r.Context(), err)
		utils.Render(w, r, envelope.SetError(err))
		return
	}
	utils.Render(w, r, envelope.SetResponse(m))
}

func delete(w http.ResponseWriter, r *http.Request) {
	envelope, ms := extract(r)
	encName := chi.URLParam(r, "name")
	name, err := url.QueryUnescape(encName)
	if err != nil {
		log.Error(r.Context(), err)
		utils.Render(w, r, envelope.SetError(err))
		return
	}
	err = ms.DeleteByName(name)
	if err != nil {
		log.Error(r.Context(), err)
		utils.Render(w, r, envelope.SetError(err))
		return
	}
	utils.Render(w, r, envelope.SetResponse(""))
}

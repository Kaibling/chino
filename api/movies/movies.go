package movies

import (
	"chino/lib/logging"
	"chino/lib/utils"
	"chino/models"
	"chino/services"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

var prepare = func(r *http.Request) (*services.MovieService, *services.CrawlerService) {
	ms := utils.GetContext("movieservice", r).(*services.MovieService)
	cs := utils.GetContext("crawlerservice", r).(*services.CrawlerService)
	return ms, cs
}

func create(w http.ResponseWriter, r *http.Request) {
	rs, _ := prepare(r)
	var model models.Movie
	err := json.NewDecoder(r.Body).Decode(&model)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	m, err := rs.Create(model)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	render.Respond(w, r, m)
}

// func readAllSaved(w http.ResponseWriter, r *http.Request) {
// 	rs, _ := prepare(r)
// 	objects, err := rs.ReadAll()
// 	if err != nil {
// 		fmt.Println(err.Error())
// 		return
// 	}
// 	render.Respond(w, r, objects)
// }

// func update(w http.ResponseWriter, r *http.Request) {
// 	w.Write([]byte("update"))
// }

func delete(w http.ResponseWriter, r *http.Request) {
	encName := chi.URLParam(r, "name")
	name, err := url.QueryUnescape(encName)
	if err != nil {
		logging.Logger.Error(err)
		w.Write([]byte(err.Error()))
		return
	}
	ms, _ := prepare(r)
	err = ms.DeleteByName(name)
	if err != nil {
		logging.Logger.Error(err)
		w.Write([]byte(err.Error()))
		return
	}
	w.Write([]byte("delete"))
}

// func readAll(w http.ResponseWriter, r *http.Request) {
// 	ms, cs := prepare(r)
// 	objects, err := cs.GetMovies(1)
// 	if err != nil {
// 		fmt.Println(err.Error())
// 		return
// 	}
// 	favorites, err := ms.ReadUntil(1)
// 	if err != nil {
// 		fmt.Println(err.Error())
// 		return
// 	}
// 	movieList := ms.FavoriteMovies(favorites, objects)
// 	render.Respond(w, r, movieList)
// }

package dashboard

import (
	"chino/lib/utils"
	"chino/services"
	"fmt"
	"net/http"
	"strconv"
	"text/template"
	"time"
)

var prepare = func(r *http.Request) (*services.MovieService, *services.CrawlerService) {
	ms := utils.GetContext("movieservice", r).(*services.MovieService)
	cs := utils.GetContext("crawlerservice", r).(*services.CrawlerService)
	return ms, cs
}

func show(w http.ResponseWriter, r *http.Request) {
	ms, cs := prepare(r)
	var month int
	key, ok := r.URL.Query()["month"]
	if !ok {
		_, m, _ := time.Now().Date()
		month = int(m)

	} else {
		m, err := strconv.Atoi(key[0])
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		month = m
	}

	objects, err := cs.GetMovies(month)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	favorites, err := ms.ReadUntil(month)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	movieList := ms.FavoriteMovies(favorites, objects)
	tmpl := template.Must(template.ParseFiles("api/dashboard/dashboard.html"))
	tmpl.Execute(w, movieList)
}

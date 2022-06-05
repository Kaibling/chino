package dashboard

import (
	"chino/lib/utils"
	"chino/services"
	"fmt"
	"net/http"
	"text/template"
)

var prepare = func(r *http.Request) (*services.MovieService, *services.CrawlerService) {
	ms := utils.GetContext("movieservice", r).(*services.MovieService)
	cs := utils.GetContext("crawlerservice", r).(*services.CrawlerService)
	return ms, cs
}

func show(w http.ResponseWriter, r *http.Request) {
	ms, cs := prepare(r)
	objects, err := cs.GetMovies(1)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	favorites, err := ms.ReadUntil(1)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	movieList := ms.FavoriteMovies(favorites, objects)
	tmpl := template.Must(template.ParseFiles("api/dashboard/dashboard.html"))
	tmpl.Execute(w, movieList)
}

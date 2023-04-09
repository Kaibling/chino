package crawler

import "github.com/Kaibling/chino/models"

type Crawler interface {
	GetMovies(months int) ([]models.Movie, error)
}

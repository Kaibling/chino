package crawler

import "chino/models"

type Crawler interface {
	GetMovies(months int) ([]models.Movie, error)
}




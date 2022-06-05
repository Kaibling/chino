package services

import (
	"chino/lib/logging"
	"chino/lib/utils"
	"chino/models"
	"errors"
	"fmt"
	"time"
)

type movieRepo interface {
	ReadAll() ([]models.Movie, error)
	ReadByID(id int) (*models.Movie, error)
	ReadByName(name string) (*models.Movie, error)
	ReadUntil(months int) ([]models.Movie, error)
	Create(models.Movie) (*models.Movie, error)
	SetNotified(name string) error
	DeleteByID(id int) error
	DeleteByName(name string) error
}

type MovieService struct {
	repo   movieRepo
	notify *NotificationService
}

func NewMovieService(repo movieRepo) *MovieService {
	return &MovieService{repo: repo}
}
func (s *MovieService) AddNotificationService(n *NotificationService) {
	s.notify = n
}

func (s *MovieService) ReadAll() ([]models.Movie, error) {
	return s.repo.ReadAll()
}
func (s *MovieService) ReadByID(id int) (*models.Movie, error) {
	return s.repo.ReadByID(id)
}
func (s *MovieService) ReadByName(name string) (*models.Movie, error) {
	return s.repo.ReadByName(name)
}

func (s *MovieService) Create(m models.Movie) (*models.Movie, error) {
	return s.repo.Create(m)
}

func (s *MovieService) DeleteByID(id int) error {
	return s.repo.DeleteByID(id)
}
func (s *MovieService) DeleteByName(name string) error {
	return s.repo.DeleteByName(name)
}

func (s *MovieService) ReadUntil(months int) ([]models.Movie, error) {
	return s.repo.ReadUntil(months)
}

func (s *MovieService) FavoriteMovies(favories []models.Movie, releases []models.Movie) []models.Movie {
	for i := range favories {
		for j := range releases {
			if favories[i].Name == releases[j].Name {
				releases[j].Favorite = true
				continue
			}
		}
	}
	return releases
}

// Scheduler function
func (s *MovieService) Run() error {
	if s.notify == nil {
		return errors.New("notification service not set")
	}
	movies, err := s.repo.ReadAll()
	if err != nil {
		logging.Logger.Error(err)
	}
	for i := range movies {
		t, err := utils.TimeFromFormat(movies[i].ReleaseDate)
		if err != nil {
			logging.Logger.Error(err)
		}
		if t.Before(time.Now().AddDate(0, 0, 20)) && !movies[i].Notified {
			s.notify.Send(fmt.Sprintf("%s will be released at %s", movies[i].Name, movies[i].ReleaseDate))
			s.repo.SetNotified(movies[i].Name)
		}
	}
	return nil
}

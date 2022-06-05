package gormrepo

import (
	"chino/lib/utils"
	"chino/models"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type Movie struct {
	DBModel
	Name        string `gorm:"unique"`
	ReleaseDate time.Time
	Genre       string
	Notified    bool
}

type MovieRepo struct {
	db *gorm.DB
}

func NewMovieRepo(db *gorm.DB) *MovieRepo {
	return &MovieRepo{db}
}
func (i *MovieRepo) ReadAll() ([]models.Movie, error) {
	movies := []Movie{}
	res := i.db.Find(&movies)
	if res.Error != nil {
		return nil, res.Error
	}
	return movieArrayUnmarshal(movies), nil
}

func (i *MovieRepo) ReadByID(id int) (*models.Movie, error) {
	return nil, nil
}
func (i *MovieRepo) SetNotified(name string) error {
	err := i.db.Table("movies").Where("name =  ?", name).Update("notified", true).Error // TODO change Table to Model
	return err
}
func (i *MovieRepo) ReadByName(name string) (*models.Movie, error) {
	movies := Movie{}
	res := i.db.Where("name = ?", name).First(&movies)
	if res.Error != nil {
		return nil, res.Error
	}
	return movieUnmarshal(movies), nil
}

func (i *MovieRepo) ReadUntil(months int) ([]models.Movie, error) {
	movies := []Movie{}
	//res := i.db.Where("release_date >= ? and release_date <= ?", time.Now(), time.Now().Add(time.Duration(months)*time.Hour*24*30)).Find(&movies)
	res := i.db.Find(&movies)
	if res.Error != nil {
		return nil, res.Error
	}
	return movieArrayUnmarshal(movies), nil
}
func (i *MovieRepo) Create(m models.Movie) (*models.Movie, error) {
	newMovie := movieMarshal(m)
	res := i.db.Create(&newMovie)
	if res.Error != nil {
		return nil, res.Error
	}
	ru := movieUnmarshal(*newMovie)
	return ru, nil
}
func (i *MovieRepo) UpdateByID(id int, m models.Movie) (*models.Movie, error) {
	panic("not implemented")
}
func (i *MovieRepo) DeleteByID(id int) error {
	panic("not implemented")
}
func (i *MovieRepo) DeleteByName(name string) error {
	res := i.db.Delete(&Movie{}, "name = ?", name)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected != 1 {
		return fmt.Errorf("delete effected not 1 row, but %d", res.RowsAffected)
	}
	return nil
}

func movieMarshal(m models.Movie) *Movie {
	t, err := utils.TimeFromFormat(m.ReleaseDate)
	if err != nil {
		panic(err)
	}
	return &Movie{
		DBModel:     DBModel{ID: m.ID},
		Name:        m.Name,
		ReleaseDate: t,
		Genre:       m.Genre,
		Notified:    m.Notified,
	}
}

func movieUnmarshal(r Movie) *models.Movie {
	return &models.Movie{
		ID:          r.ID,
		Name:        r.Name,
		ReleaseDate: utils.TimeToFormat(r.ReleaseDate),
		Genre:       r.Genre,
		Favorite:    true,
		Notified:    r.Notified,
	}
}

// func movieArrayMarshal(m []models.Movie) []Movie {
// 	ms := []Movie{}
// 	for i := range m {
// 		ms = append(ms, *movieMarshal(m[i]))
// 	}
// 	return ms
// }

func movieArrayUnmarshal(m []Movie) []models.Movie {
	ms := []models.Movie{}
	for i := range m {
		ms = append(ms, *movieUnmarshal(m[i]))
	}
	return ms
}

package repo_sqlx

import (
	"context"
	"fmt"
	"time"

	"github.com/Kaibling/chino/models"
	"github.com/Kaibling/chino/pkg/utils"

	"github.com/jmoiron/sqlx"
)

type movie struct {
	dbModel
	Name        string    `db:"name"`
	ReleaseDate time.Time `db:"release_date"`
	Genre       string    `db:"genre"`
	Notified    bool      `db:"notified"`
}

type MovieRepo struct {
	db  *sqlx.DB
	ctx context.Context
	// loggerFields []log.Field
}

func NewMovieRepo(ctx context.Context, db *sqlx.DB) *MovieRepo {
	return &MovieRepo{db: db, ctx: ctx}
}

func (r *MovieRepo) ReadAll() ([]models.Movie, error) {
	movies := []movie{}
	if err := r.db.Select(&movies, "SELECT * FROM movies"); err != nil {
		return nil, err
	}
	return movieArrayUnmarshal(movies), nil
}

func (r *MovieRepo) ReadByID(id int) (*models.Movie, error) {
	return nil, nil
}

func (r *MovieRepo) SetNotified(name string) error {
	res, err := r.db.Exec(`UPDATE movies SET notified = ? WHERE name = ?;`, true, name)
	if err != nil {
		return err
	}
	re, _ := res.RowsAffected()
	if re == 0 {
		return fmt.Errorf("movie notifiaction could not be set to true")
	}
	return err
}

func (r *MovieRepo) ReadByName(name string) (*models.Movie, error) {
	movie := movie{}
	if err := r.db.Select(&movie, "SELECT * FROM movies where name = ?;", name); err != nil {
		return nil, err
	}
	return movieUnmarshal(movie), nil
}

func (r *MovieRepo) ReadUntil(months int) ([]models.Movie, error) {
	movies := []movie{}
	if err := r.db.Select(&movies, "SELECT * FROM movies where release_date >= ? and release_date <= ?;", time.Now(), time.Now().Add(time.Duration(months)*time.Hour*24*30)); err != nil {
		return nil, err
	}
	return movieArrayUnmarshal(movies), nil
}

func (r *MovieRepo) Create(newMovie models.Movie) (*models.Movie, error) {
	m := movieMarshal(newMovie)
	m.CreatedAt = time.Now()
	//m.CreatedBy = r.ctx.Value("user").(string)
	m.ID = utils.NewULID().String()
	query := `INSERT INTO movies(name,release_date, genre,notified,id,created_at,created_by) 
          VALUES(:name, :release_date,:genre,:notified,:id,:created_at,:created_by)`
	res, err := r.db.NamedExec(query, m)
	if err != nil {
		return nil, err
	}
	re, _ := res.RowsAffected()
	if re == 0 {
		return nil, fmt.Errorf("no movie inserted")
	}
	return movieUnmarshal(*m), nil
}

func (r *MovieRepo) DeleteByName(name string) error {
	if res, err := r.db.Exec(`delete from movies where name = ?;`, name); err != nil {
		return err
	} else {
		if re, _ := res.RowsAffected(); re > 1 {
			return fmt.Errorf("delete effected not 1 row, but %d", re)
		} else if re, _ := res.RowsAffected(); re == 0 {
			return fmt.Errorf("no rows deleted")
		}
	}
	return nil
}

func movieMarshal(m models.Movie) *movie {
	t, err := utils.TimeFromFormat(m.ReleaseDate)
	if err != nil {
		panic(err)
	}
	return &movie{
		dbModel:     dbModel{ID: m.ID},
		Name:        m.Name,
		ReleaseDate: t,
		Genre:       m.Genre,
		Notified:    m.Notified,
	}
}

func movieUnmarshal(r movie) *models.Movie {
	return &models.Movie{
		ID:          r.ID,
		Name:        r.Name,
		ReleaseDate: utils.TimeToFormat(r.ReleaseDate),
		Genre:       r.Genre,
		Favorite:    true,
		Notified:    r.Notified,
	}
}

func movieArrayUnmarshal(m []movie) []models.Movie {
	ms := []models.Movie{}
	for i := range m {
		ms = append(ms, *movieUnmarshal(m[i]))
	}
	return ms
}

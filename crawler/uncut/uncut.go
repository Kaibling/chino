package uncut

import (
	"bytes"
	"context"
	"strconv"
	"strings"
	"time"

	"chino/models"
	"chino/pkg/log"
	"chino/pkg/utils"

	"github.com/PuerkitoBio/goquery"
)

type Uncut struct {
	ctx context.Context
}

func NewCrawler(ctx context.Context) *Uncut {
	return &Uncut{ctx: ctx}
}

func (s *Uncut) GetMovies(months int) ([]models.Movie, error) {
	year, _, _ := time.Now().AddDate(0, months, 0).Date()
	currentYear, _, _ := time.Now().Date()
	diff := year - currentYear
	movies := []models.Movie{}
	for i := 0; i <= diff; i++ {
		y := strconv.Itoa(currentYear + i)
		moviesPerYear, err := getYear(s.ctx, y, months)
		if err != nil {
			return nil, err
		}
		movies = append(movies, moviesPerYear...)
	}
	return movies, nil
}

func timeParse(date string) (time.Time, error) {
	layout := "02.01.2006"
	date = strings.ReplaceAll(date, " ", "")
	t, err := time.Parse(layout, date)
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}

func getYear(ctx context.Context, y string, months int) ([]models.Movie, error) {
	url := "https://www.uncut.at/movies/jahr.php?country=AT&year=" + y
	movies := []models.Movie{}
	log.Info(ctx, url)
	result, err := utils.Request("GET", url, nil)
	if err != nil {
		return nil, err
	}
	resultBuffer := bytes.NewBuffer(result)
	doc, err := goquery.NewDocumentFromReader(resultBuffer)
	if err != nil {
		log.Error(ctx, err)
	}
	tab := doc.Find(".tabelle1")
	tab.Find("tr").Each(func(i int, s *goquery.Selection) {
		if len(s.Children().Nodes) == 4 {
			movie := models.Movie{}
			for i := range s.Children().Nodes {

				c := s.Children().Nodes[i]
				switch i {
				case 0:
					rd, err := timeParse(c.FirstChild.Data + y)
					if err != nil {
						log.Error(ctx, err)
						continue
					}
					if rd.Before(time.Now()) || rd.After(time.Now().AddDate(0, months, 0)) {
						continue
					}
					movie.ReleaseDate = utils.TimeToFormat(rd)
				case 1:
					movie.Name = c.FirstChild.FirstChild.FirstChild.Data
				case 2:
					movie.Genre = c.FirstChild.Data
				}
			}
			if movie.ReleaseDate != "" {
				movies = append(movies, movie)
			}
		}
	})
	return movies, nil
}

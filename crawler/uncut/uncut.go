package uncut

import (
	"bytes"
	"chino/lib/logging"
	"chino/lib/utils"
	"chino/models"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type uncut struct {
}

func NewCrawler() *uncut {
	return &uncut{}
}

func (s *uncut) GetMovies(months int) ([]models.Movie, error) {
	year, _, _ := time.Now().AddDate(0, months, 0).Date()
	currentYear, _, _ := time.Now().Date()
	diff := year - currentYear
	movies := []models.Movie{}
	for i := 0; i <= diff; i++ {
		y := strconv.Itoa(currentYear + i)
		moviesPerYear, err := getYear(y, months)
		if err != nil {
			return nil, err
		}
		movies = append(movies, moviesPerYear...)
	}
	return movies, nil
}

func timeParse(date string) time.Time {

	layout := "02.01.2006"
	date = strings.ReplaceAll(date, " ", "")
	t, err := time.Parse(layout, date)
	if err != nil {
		fmt.Println(err)
	}
	return t
}

func getYear(y string, months int) ([]models.Movie, error) {
	url := "https://www.uncut.at/movies/jahr.php?country=AT&year=" + y
	movies := []models.Movie{}
	logging.Logger.Info(url)
	result, err := utils.Request("GET", url, nil)
	if err != nil {
		return nil, err
	}
	resultBuffer := bytes.NewBuffer(result)
	doc, err := goquery.NewDocumentFromReader(resultBuffer)
	if err != nil {
		log.Fatal(err)
	}
	tab := doc.Find(".tabelle1")
	tab.Find("tr").Each(func(i int, s *goquery.Selection) {
		if len(s.Children().Nodes) == 4 {
			movie := models.Movie{}
			for i := range s.Children().Nodes {

				c := s.Children().Nodes[i]
				switch i {
				case 0:
					rd := timeParse(c.FirstChild.Data + y)
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

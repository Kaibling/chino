package uncut

import (
	"bytes"
	"chino/lib/utils"
	"chino/models"
	"fmt"
	"log"
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
	movies := []models.Movie{}
	year := "2022"
	// TODO months to years
	url := "https://www.uncut.at/movies/jahr.php?country=AT&year=" + year

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
					rd := timeParse(c.FirstChild.Data + year)
					if rd.Before(time.Now()) {
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

func timeParse(date string) time.Time {

	layout := "02.01.2006"
	date = strings.ReplaceAll(date, " ", "")
	t, err := time.Parse(layout, date)
	if err != nil {
		fmt.Println(err)
	}
	return t
}

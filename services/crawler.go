package services

import (
	"github.com/Kaibling/chino/crawler"
	"github.com/Kaibling/chino/models"
)

type CrawlerService struct {
	c crawler.Crawler
}

func NewCrawlerService(c crawler.Crawler) *CrawlerService {
	return &CrawlerService{c}
}

func (s *CrawlerService) GetMovies(m int) ([]models.Movie, error) {
	return s.c.GetMovies(m)
}

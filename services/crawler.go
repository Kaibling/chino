package services

import (
	"chino/crawler"
	"chino/models"
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

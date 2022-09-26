package services

import (
	"lihood/g"
	"lihood/internal/repository"
)

type RecommendService interface {
	Count() (count int64, err error)
}

func NewRecommendService(session g.Session) RecommendService {
	return &recommendService{session: session}
}

type recommendService struct {
	session g.Session
}

func (s recommendService) Count() (count int64, err error) {
	return repository.NewRecommendRepository(s.session).Count()
}

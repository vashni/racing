package service

import (
	"github.com/vashni/racing/sports/db"
	"github.com/vashni/racing/sports/proto/sports"
	"golang.org/x/net/context"
)

type Sports interface {
	// ListSports  will return a collection of sport events.
	ListEvents(ctx context.Context, in *sports.ListEventsRequest) (*sports.ListEventsResponse, error)
}

// racingService implements the Racing interface.
type sportsService struct {
	sportsRepo db.SportsRepo
}

// NewSportsService instantiates and returns a new racingService.
func NewSportsService(sportsRepo db.SportsRepo) Sports {
	return &sportsService{sportsRepo}
}

func (s *sportsService) ListEvents(ctx context.Context, in *sports.ListEventsRequest) (*sports.ListEventsResponse, error) {
	eventsRecords, err := s.sportsRepo.List(in, in.OrderBy)
	if err != nil {
		return nil, err
	}
	return &sports.ListEventsResponse{Event: eventsRecords}, nil
}

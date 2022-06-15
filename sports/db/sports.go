package db

import (
	"database/sql"
	"github.com/golang/protobuf/ptypes"
	_ "github.com/mattn/go-sqlite3"
	"sync"
	"time"

	"github.com/vashni/racing/sports/proto/sports"
)

// SportRepo provides repository access to sport related events and other data.
type SportsRepo interface {
	// Init will initialise our sports repository.
	Init() error

	// List will return a list of sport.
	List(eventRequest *sports.ListEventsRequest, orderBy string) ([]*sports.Event, error)
	//Get(raceIdRequest *sports.GetRaceRequest) (*sports.Race, error)
}

type sportsRepo struct {
	db   *sql.DB
	init sync.Once
}

// NewSportsRepo creates a new sports repository.
func NewSportsRepo(db *sql.DB) SportsRepo {
	return &sportsRepo{db: db}
}

// Init prepares the sport repository dummy data.
func (r *sportsRepo) Init() error {
	var err error

	r.init.Do(func() {
		// For test/example purposes, we seed the DB with some dummy sport.
		err = r.seed()
	})

	return err
}

func (r *sportsRepo) List(filter *sports.ListEventsRequest, orderBy string) ([]*sports.Event, error) {
	var (
		err   error
		query string
		args  []interface{}
	)

	query = getSportsQueries()[sportsList]

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}

	return r.scanEvents(rows)
}

func (m *sportsRepo) scanEvents(
	rows *sql.Rows,
) ([]*sports.Event, error) {
	var eventRecords []*sports.Event

	for rows.Next() {
		var event sports.Event
		var advertisedStart time.Time

		if err := rows.Scan(&event.Id, &event.Name, &advertisedStart); err != nil {
			if err == sql.ErrNoRows {
				return nil, nil
			}

			return nil, err
		}

		ts, err := ptypes.TimestampProto(advertisedStart)
		if err != nil {
			return nil, err
		}

		event.AdvertisedStartTime = ts
		eventRecords = append(eventRecords, &event)
	}

	return eventRecords, nil
}

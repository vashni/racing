package db

import (
	"database/sql"
	"errors"
	"github.com/golang/protobuf/ptypes"
	_ "github.com/mattn/go-sqlite3"
	"strconv"
	"strings"
	"sync"
	"time"

	"git.neds.sh/matty/entain/racing/proto/racing"
)

// RacesRepo provides repository access to races.
type RacesRepo interface {
	// Init will initialise our races repository.
	Init() error

	// List will return a list of races.
	List(filter *racing.ListRacesRequestFilter, orderBy string) ([]*racing.Race, error)
}

type racesRepo struct {
	db   *sql.DB
	init sync.Once
}

// NewRacesRepo creates a new races repository.
func NewRacesRepo(db *sql.DB) RacesRepo {
	return &racesRepo{db: db}
}

// Init prepares the race repository dummy data.
func (r *racesRepo) Init() error {
	var err error

	r.init.Do(func() {
		// For test/example purposes, we seed the DB with some dummy races.
		err = r.seed()
	})

	return err
}

func (r *racesRepo) List(filter *racing.ListRacesRequestFilter, orderBy string) ([]*racing.Race, error) {
	var (
		err   error
		query string
		args  []interface{}
	)

	query = getRaceQueries()[racesList]

	query, args = r.applyFilter(query, filter)

	query, err = r.applyOrderBy(query, orderBy) //TODO: Need to handle more error cases.

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}

	return r.scanRaces(rows)
}

//Expecting order by in the format columnName.order
func (ob *racesRepo) applyOrderBy(query string, orderBy string) (string, error) {
	if orderBy == "" {
		// Default sort column will be advertised_start_time.
		orderBy = "advertised_start_time.asc"
	}

	splits := strings.Split(orderBy, ".")
	if len(splits) != 2 {
		return "", errors.New("proper format order by needed column name order")
	}

	query = query + " order by " + splits[0] + " " + splits[1]

	return query, nil
}

func (r *racesRepo) applyFilter(query string, filter *racing.ListRacesRequestFilter) (string, []interface{}) {
	var (
		clauses []string
		args    []interface{}
	)
	if filter == nil {
		return query, args
	}

	if len(filter.MeetingIds) > 0 {
		clauses = append(clauses, "meeting_id IN ("+strings.Repeat("?,", len(filter.MeetingIds)-1)+"?)")

		for _, meetingID := range filter.MeetingIds {
			args = append(args, meetingID)
		}
	}

	//filter to fetch the visible races.
	if filter.IsVisible {
		clauses = append(clauses, " visible = "+strconv.FormatBool(filter.IsVisible))
	}

	if len(clauses) != 0 {
		query += " WHERE " + strings.Join(clauses, " AND ")
	}

	return query, args
}

func (m *racesRepo) scanRaces(
	rows *sql.Rows,
) ([]*racing.Race, error) {
	var races []*racing.Race

	for rows.Next() {
		var race racing.Race
		var advertisedStart time.Time

		if err := rows.Scan(&race.Id, &race.MeetingId, &race.Name, &race.Number, &race.Visible, &advertisedStart); err != nil {
			if err == sql.ErrNoRows {
				return nil, nil
			}

			return nil, err
		}

		ts, err := ptypes.TimestampProto(advertisedStart)
		if err != nil {
			return nil, err
		}

		race.AdvertisedStartTime = ts

		//Task3 - set status closed/open based on the past advertised_start_date time
		dt := time.Now()
		if dt.After(advertisedStart) {
			race.Status = "CLOSED"
		} else {
			race.Status = "OPEN"
		}
		races = append(races, &race)
	}

	return races, nil
}

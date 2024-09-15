package db

import (
	"fmt"
	"sync"
	"time"

	"github.com/upper/db/v4"
	"github.com/upper/db/v4/adapter/cockroachdb"
)

var dbSess db.Session
var mutex sync.Mutex

type InstrumentInfo struct {
	ID               int64     `db:"id,omitempty"`
	UpdateTime       time.Time `db:"update_time"`
	Instrument       string    `db:"instrument"`
	Price            float64   `db:"price"`
	QuotedInstrument string    `db:"quoted_instrument"`
}

func OpenDB() error {
	settings := cockroachdb.ConnectionURL{
		Database: `trading`,
		Host:     `localhost`,
		User:     `root`,
		Password: ``,
		Options: map[string]string{
			"sslmode": "disable",
		},
	}

	var err error
	dbSess, err = cockroachdb.Open(settings)
	if err != nil {
		return fmt.Errorf("database connection error: %w", err)
	}

	return nil
}

func CloseDB() {
	if dbSess != nil {
		dbSess.Close()
	}
}

func SaveInstrumentInfoToDB(info InstrumentInfo) error {
	if dbSess == nil {
		return fmt.Errorf("database session is not initialized")
	}

	mutex.Lock()
	defer mutex.Unlock()

	_, err := dbSess.Collection("prices").
		Insert(info)
	if err != nil {
		return fmt.Errorf("saving to database error: %w", err)
	}

	return nil
}

func GetInstrumentInfoFromDB(instrument string, quotedInstrument string) (InstrumentInfo, error) {
	if dbSess == nil {
		return InstrumentInfo{}, fmt.Errorf("database session is not initialized")
	}

	mutex.Lock()
	defer mutex.Unlock()

	var info InstrumentInfo
	err := dbSess.Collection("prices").
		Find(db.Cond{"instrument": instrument, "quoted_instrument": quotedInstrument}).
		OrderBy("-update_time").
		One(&info)

	if err != nil {
		return InstrumentInfo{}, fmt.Errorf("getting data from database error: %w", err)
	}

	return info, nil
}
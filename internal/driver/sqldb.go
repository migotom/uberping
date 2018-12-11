package driver

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/lib/pq" // load psql driver
	"github.com/migotom/uberping/internal/schema"
)

type sqlDB struct {
	conn     *sql.DB
	dbConfig *schema.DBConfig
}

func (d *sqlDB) connect() error {
	var err error
	d.conn, err = sql.Open(d.dbConfig.Driver, d.dbConfig.Params)
	if err != nil {
		return err
	}
	return nil
}

type retryFunc func() error

func (d *sqlDB) retry(retryFunc retryFunc) (err error) {
	for retries := 0; retries < 3; retries++ {
		err = retryFunc()
		if err != nil {
			// cleanup
			d.conn.Close()

			// reconnect and retry
			time.Sleep(1000 * time.Millisecond)
			d.connect()
			continue
		}
	}
	return
}

func (d *sqlDB) Query(query string, args ...interface{}) (rows *sql.Rows, err error) {
	err = d.retry(func() error {
		rows, err = d.conn.Query(query, args...)
		return err
	})
	return
}

func (d *sqlDB) Exec(query string, args ...interface{}) (result sql.Result, err error) {
	err = d.retry(func() error {
		result, err = d.conn.Exec(query, args...)
		return err
	})
	return
}

func getDB(dbConfig *schema.DBConfig) *sqlDB {
	db, ok := dbConfig.Connection.(*sqlDB)
	if !ok {
		db = &sqlDB{}
		db.dbConfig = dbConfig
		dbConfig.Connection = db
	}
	return db
}

// DBCleaner closes DB connection.
func DBCleaner(dbConfig *schema.DBConfig) {
	db, ok := dbConfig.Connection.(*sqlDB)
	if ok {
		defer db.conn.Close()
	}
}

// DBSqlLoadHosts loads list of hosts from database.
func DBSqlLoadHosts(hostParser schema.HostParser, dbConfig *schema.DBConfig) ([]schema.Host, error) {
	db := getDB(dbConfig)
	if err := db.connect(); err != nil {
		return nil, err
	}

	var hosts []schema.Host

	rows, err := db.Query(dbConfig.Queries.GetDevices, dbConfig.IDserver)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		host := schema.Host{}

		err = rows.Scan(&host.ID, &host.IP, &host.InactiveSince)
		if err != nil {
			return nil, err
		}
		host.IP, host.Port, err = hostParser(host.IP)
		if err != nil {
			return nil, err
		}

		hosts = append(hosts, host)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return hosts, nil
}

// DBSqlSavePingResult save ping results using db.
func DBSqlSavePingResult(result schema.ProbeResult, dbConfig *schema.DBConfig) error {
	db, ok := dbConfig.Connection.(*sqlDB)
	if !ok {
		log.Fatal("No database connection")
	}

	if result.Loss == 100 {
		if !result.Host.InactiveSince.Valid {
			result.Host.InactiveSince = sql.NullString{String: "NOW()", Valid: true}
		}
	} else {
		result.Host.InactiveSince = sql.NullString{}
	}

	_, err := db.Exec(dbConfig.Queries.UpdateDevice, result.Loss, result.AvgTime, result.Host.InactiveSince, result.Host.ID)
	if err != nil {
		return err
	}

	return nil
}

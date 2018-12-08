package driver

import (
	"database/sql"
	"log"

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
	//defer db.Close()
	return nil
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

	rows, err := db.conn.Query(dbConfig.Queries.GetDevices, dbConfig.IDserver)
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
		host.IP, err = hostParser(host.IP)
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
func DBSqlSavePingResult(result schema.PingResult, dbConfig *schema.DBConfig) error {
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

	_, err := db.conn.Exec(dbConfig.Queries.UpdateDevice, result.Loss, result.AvgTime, result.Host.InactiveSince, result.Host.ID)
	if err != nil {
		return err
	}

	return nil
}

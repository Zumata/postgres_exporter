package stats

import (
	"database/sql"

	"github.com/Zumata/exporttools"
)

type DatabaseCollector struct {
	db     *sql.DB
	dbHost string
	dbName string
}

func NewDatabaseCollector(db *sql.DB, dbHost, dbName string) *DatabaseCollector {
	return &DatabaseCollector{db, dbHost, dbName}
}

func (s *DatabaseCollector) Collect() ([]*exporttools.Metric, error) {
	stat := new(DatabaseStat)
	err := collectDatabaseStats(s.db, s.dbName, stat)
	if err != nil {
		return make([]*exporttools.Metric, 0), err
	}
	return formatDatabaseStats(s.dbHost, stat), nil
}

type DatabaseStat struct {
	Name                                          string
	Commit, Rollback                              int64
	Read, Hit                                     int64
	Returned, Fetched, Inserted, Updated, Deleted int64
	Conflicts, TempFiles, TempBytes, Deadlocks    int64
	ReadTime, WriteTime                           float64
}

func collectDatabaseStats(db *sql.DB, dbName string, s *DatabaseStat) error {
	var ignore interface{}
	return db.QueryRow(`SELECT * FROM pg_stat_database WHERE datname=$1`, dbName).Scan(
		&ignore, &s.Name, &ignore,
		&s.Commit, &s.Rollback,
		&s.Read, &s.Hit,
		&s.Returned, &s.Fetched, &s.Inserted, &s.Updated, &s.Deleted,
		&s.Conflicts, &s.TempFiles, &s.TempBytes, &s.Deadlocks,
		&s.ReadTime, &s.WriteTime,
		&ignore,
	)
}

func formatDatabaseStats(dbHost string, s *DatabaseStat) []*exporttools.Metric {

	metrics := []*exporttools.Metric{
		{
			Name:        "commits",
			Type:        exporttools.Gauge,
			Value:       s.Commit,
			Description: "commits",
		},
		{
			Name:        "rollbacks",
			Type:        exporttools.Gauge,
			Value:       s.Rollback,
			Description: "rollbacks",
		},
		{
			Name:        "reads",
			Type:        exporttools.Gauge,
			Value:       s.Read,
			Description: "reads",
		},
		{
			Name:        "hits",
			Type:        exporttools.Gauge,
			Value:       s.Hit,
			Description: "hits",
		},
		{
			Name:        "returns",
			Type:        exporttools.Gauge,
			Value:       s.Commit,
			Description: "returns",
		},
		{
			Name:        "fetches",
			Type:        exporttools.Gauge,
			Value:       s.Fetched,
			Description: "fetches",
		},
		{
			Name:        "inserts",
			Type:        exporttools.Gauge,
			Value:       s.Inserted,
			Description: "inserts",
		},
		{
			Name:        "updates",
			Type:        exporttools.Gauge,
			Value:       s.Updated,
			Description: "updates",
		},
		{
			Name:        "deletes",
			Type:        exporttools.Gauge,
			Value:       s.Deleted,
			Description: "deletes",
		},
		{
			Name:        "conflicts",
			Type:        exporttools.Gauge,
			Value:       s.Conflicts,
			Description: "conflicts",
		},
		{
			Name:        "temp_files",
			Type:        exporttools.Gauge,
			Value:       s.TempFiles,
			Description: "temp_files",
		},
		{
			Name:        "temp_bytes",
			Type:        exporttools.Gauge,
			Value:       s.TempBytes,
			Description: "temp_bytes",
		},
		{
			Name:        "deadlocks",
			Type:        exporttools.Gauge,
			Value:       s.Commit,
			Description: "deadlocks",
		},
	}

	for idx := range metrics {
		metrics[idx].LabelKeys = []string{"addr", "database"}
		metrics[idx].LabelVals = []string{dbHost, s.Name}
	}

	return metrics

}

package postgres

import (
	"database/sql"
	"net/url"

	"github.com/Zumata/exporttools"
	"github.com/Zumata/postgres_exporter/postgres/stats"
	_ "github.com/lib/pq"

	"github.com/prometheus/client_golang/prometheus"
)

type postgresExporter struct {
	*exporttools.BaseExporter
	connURL string
	dbName  string
	db      *sql.DB
}

func NewExporter(connURL, dbName string) *postgresExporter {
	e := &postgresExporter{
		BaseExporter: exporttools.NewBaseExporter("postgres"),
		connURL:      connURL,
		dbName:       dbName,
	}
	return e
}

func (e *postgresExporter) Setup() error {
	db, err := sql.Open("postgres", e.connURL)
	if err != nil {
		return err
	}
	parsedURL, err := url.Parse(e.connURL)
	if err != nil {
		return err
	}
	e.AddGroup(stats.NewDatabaseCollector(db, parsedURL.Host, e.dbName))
	return nil
}

func (e *postgresExporter) Close() error {
	if e.db != nil {
		return e.db.Close()
	}
	return nil
}

func (e *postgresExporter) Describe(ch chan<- *prometheus.Desc) {
	exporttools.GenericDescribe(e.BaseExporter, ch)
}

func (e *postgresExporter) Collect(ch chan<- prometheus.Metric) {
	exporttools.GenericCollect(e.BaseExporter, ch)
}

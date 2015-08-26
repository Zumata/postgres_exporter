package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/Zumata/exporttools"
	"github.com/Zumata/postgres_exporter/postgres"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	listenAddress = flag.String("web.listen-address", ":9132", "Address to listen on for web interface and telemetry.")
	metricPath    = flag.String("web.telemetry-path", "/metrics", "Path under which to expose metrics.")
	connURL       = flag.String("db.url", "postgres://user:password@localhost/dbname?sslmode=disable", "DB Url")
	dbName        = flag.String("db.name", "dbname", "DB Name")
)

func main() {

	flag.Parse()

	exporter := postgres.NewExporter(*connURL, *dbName)
	err := exporttools.Export(exporter)
	if err != nil {
		log.Fatal(err)
	}

	http.Handle(*metricPath, prometheus.Handler())
	http.HandleFunc("/", exporttools.DefaultMetricsHandler("Postgres exporter", *metricPath))
	err = http.ListenAndServe(*listenAddress, nil)
	if err != nil {
		log.Fatal(err)
	}

}

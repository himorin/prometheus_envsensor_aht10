package main

import (
  "flag"
  "log"
  "net/http"

	"github.com/prometheus/client_golang/prometheus"
  "github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
  namespace = "envsensor"
)


var (
	port = flag.String("port", "9901", "The port to listen as prometheus node")
	metricsPath = flag.String("prometheus path", "/metrics", "URL path for collected metrics")

	config = flag.String("config", "/etc/default/prometheus_envsensor", "File name of configuration")
)

func main() {
	flag.Parse()
	list_addr := "0.0.0.0:" + *port

	prometheus.MustRegister(NewAHT10Exporter(2, 0x38))
	
	http.Handle(*metricsPath, promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, *metricsPath, http.StatusMovedPermanently)
	})
	log.Fatal(http.ListenAndServe(list_addr, nil))
}


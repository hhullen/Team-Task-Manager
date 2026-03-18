package metrics

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	address      = ":2112"
	readTimeout  = time.Second * 5
	writeTimeout = time.Second * 5
	metricsPath  = "/metrics"
)

var (
	responses = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "team_task_manager_responses",
		Help:    "team_task_manager_responses",
		Buckets: []float64{0.003, 0.01, 0.05, 0.1, 0.2, 0.5, 1, 2, 5},
	},
		[]string{"method", "url", "status_code"},
	)

	dbStats = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "team_task_manager_db_state",
		Help: "team_task_manager_db_state",
	},
		[]string{"state"})
)

func init() {
	metricsMux := http.NewServeMux()
	metricsMux.Handle(metricsPath, promhttp.Handler())

	server := &http.Server{
		Addr:         address,
		Handler:      metricsMux,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
	}

	go func() {
		_ = server.ListenAndServe()
	}()
}

func ReportResponse(method, url string, statusCode int, miliseconds int64) error {
	obs, err := responses.GetMetricWithLabelValues(method, url, strconv.Itoa(statusCode))
	if err != nil {
		return err
	}

	obs.Observe(float64(miliseconds) / 1000.0)

	return nil
}

func RepordDBStats(stats sql.DBStats) {
	dbStats.WithLabelValues("in_use").Set(float64(stats.InUse))
	dbStats.WithLabelValues("idle").Set(float64(stats.Idle))
	dbStats.WithLabelValues("wait_count").Set(float64(stats.WaitCount))
	dbStats.WithLabelValues("wait_duration").Set(float64(stats.WaitDuration))
	dbStats.WithLabelValues("max_open_conns").Set(float64(stats.MaxOpenConnections))
}

package metrics

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	responseLatency = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: "response_latency",
		Help: "latency of the response",
	},
		// favor adicionar itens à essa lista com cautela.
		// uma alta cardinalidade em vetores de uma métrica no prometheus
		// pode causar impactos na performance do servidor remoto.
		[]string{"status", "method", "path", "panic"},
	)

	TransactionLatency = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: "transaction_db_latency",
		Help: "latency of the response",
	}, []string{})
)

func ResponseLatencyObserve(latency time.Duration, statusCode int, panicked bool, r *http.Request) {
	var (
		method = r.Method
		path   = r.URL.Path
	)
	observer, err := responseLatency.GetMetricWithLabelValues(
		strconv.Itoa(statusCode),
		method,
		path,
		fmt.Sprint(panicked),
	)
	if err != nil {
		// não temos como tratar esse erro
		panic(err)
	}
	observer.Observe(latency.Seconds())
}

package log

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	log "github.com/sirupsen/logrus"
)

type Context struct {
	ResponseStatus int    `json:"responseStatus"`
	Operation      string `json:"operation"`
	Service        string `json:"service"`
	ResponseTime   int64  `json:"responseTime"`
}

type Counters struct {
	HTTPRequestsResponseTime *prometheus.GaugeVec
	TotalRequests            prometheus.Counter
	HTTPRequests             *prometheus.CounterVec
}

func SetupCounters() *Counters {
	var counters Counters
	counters.TotalRequests = promauto.NewCounter(prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "The total number of requests served",
	})
	counters.HTTPRequests = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "api_http_request_status",
		Help: "http request breakdown",
	}, []string{"service", "status", "operation"})
	counters.HTTPRequestsResponseTime = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "api_avg_response_time",
		Help: "http request response time",
	}, []string{"service", "operation"})
	return &counters
}

func InstrumentedHandler(fn func(http.ResponseWriter, *http.Request), counters *Counters) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		recorder := newMetricsRecorder(w, counters)
		fn(&recorder, r)
		recorder.Replay()
	}
}

func RecordMetrics(op string, svc string, w http.ResponseWriter) {
	var ctx Context
	ctx.Operation = op
	ctx.Service = svc
	ctx.ResponseTime = time.Now().UnixNano()
	bytes, err := json.Marshal(ctx)
	if err != nil {
		log.Errorf("Marshal error %s", err)
		w.WriteHeader(500)
		return
	}
	_, err = w.Write(bytes)
	if err != nil {
		log.Errorf("Write error recording metrics %s", err)
		w.WriteHeader(500)
	}
}

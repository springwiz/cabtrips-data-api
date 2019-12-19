package log

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

type MetricsRecorder struct {
	Context
	writer   http.ResponseWriter
	counters *Counters
}

func (m *MetricsRecorder) Header() http.Header {
	return m.writer.Header()
}

func (m *MetricsRecorder) Write(bytes []byte) (int, error) {
	var ctx Context
	err := json.Unmarshal(bytes, &ctx)
	if err != nil {
		return m.writer.Write(bytes)
	} else {
		m.Context = ctx
		return len(bytes), nil
	}
}

func (m *MetricsRecorder) WriteHeader(statusCode int) {
	m.writer.WriteHeader(statusCode)
	m.Context.ResponseTime = time.Now().UnixNano() - m.Context.ResponseTime
	m.Context.ResponseStatus = statusCode
}

func newMetricsRecorder(writer http.ResponseWriter, counters *Counters) MetricsRecorder {
	return MetricsRecorder{
		Context:  Context{},
		writer:   writer,
		counters: counters,
	}
}

func (m *MetricsRecorder) Replay() {
	m.counters.TotalRequests.Inc()
	m.counters.HTTPRequests.WithLabelValues(m.Context.Service, m.Context.Operation,
		strconv.Itoa(m.Context.ResponseStatus)).Inc()
	m.counters.HTTPRequestsResponseTime.WithLabelValues(m.Context.Service,
		m.Context.Operation).Set(float64(m.Context.ResponseTime))
}

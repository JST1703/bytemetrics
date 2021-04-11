package bytesize

import (
	prom "github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc/stats"
)

type ServerByteMetrics struct {
	serverMsgSizeReceivedHistogramOpts prom.HistogramOpts
	serverMsgSizeReceivedHistogram     *prom.HistogramVec

	serverMsgSizeSentHistogramOpts prom.HistogramOpts
	serverMsgSizeSentHistogram     *prom.HistogramVec
}

func NewServerByteMetrics() *ServerByteMetrics {
	tempreceived := prom.HistogramOpts{
		Name:    "grpc_server_msg_size_received_bytes",
		Help:    "Histogram of message sizes received by the server.",
		Buckets: defMsgBytesBuckets,
	}
	tempsent := prom.HistogramOpts{
		Name:    "grpc_server_msg_size_sent_bytes",
		Help:    "Histogram of message sizes sent by the server.",
		Buckets: defMsgBytesBuckets,
	}
	return &ServerByteMetrics{
		serverMsgSizeReceivedHistogramOpts: tempreceived,
		serverMsgSizeReceivedHistogram: prom.NewHistogramVec(
			tempreceived,
			[]string{"grpc_service", "grpc_method", "grpc_stats"},
		),
		serverMsgSizeSentHistogramOpts: tempsent,
		serverMsgSizeSentHistogram: prom.NewHistogramVec(
			tempsent,
			[]string{"grpc_service", "grpc_method", "grpc_stats"},
		),
	}
}

// Describe sends the super-set of all possible descriptors of metrics
// collected by this Collector to the provided channel and returns once
// the last descriptor has been sent.
func (m *ServerByteMetrics) Describe(ch chan<- *prom.Desc) {
	m.serverMsgSizeReceivedHistogram.Describe(ch)
	m.serverMsgSizeSentHistogram.Describe(ch)
}

// Collect is called by the Prometheus registry when collecting
// metrics. The implementation sends each collected metric via the
// provided channel and returns once the last metric has been sent.
func (m *ServerByteMetrics) Collect(ch chan<- prom.Metric) {
	m.serverMsgSizeReceivedHistogram.Collect(ch)
	m.serverMsgSizeSentHistogram.Collect(ch)
}

func (m *ServerByteMetrics) NewServerByteStatsHandler() stats.Handler {
	return &serverByteStatsHandler{
		serverByteMetrics: m,
	}
}

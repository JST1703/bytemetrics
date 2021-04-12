package bytesize

import (
	prom "github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc/stats"
)

type ServerByteMetrics struct {
	serverMsgSizeBytesSent     *prom.CounterVec
	serverMsgSizeBytesReceived *prom.CounterVec

	serverMsgSizeReceivedHistogramEnabled bool
	serverMsgSizeReceivedHistogramOpts    prom.HistogramOpts
	serverMsgSizeReceivedHistogram        *prom.HistogramVec

	serverMsgSizeSentHistogramEnabled bool
	serverMsgSizeSentHistogramOpts    prom.HistogramOpts
	serverMsgSizeSentHistogram        *prom.HistogramVec
}

func NewServerByteMetrics() *ServerByteMetrics {
	return &ServerByteMetrics{
		serverMsgSizeBytesReceived: prom.NewCounterVec(
			prom.CounterOpts{
				Name: "grpc_server_bytes",
				Help: "Number of message sizes received by the server.",
			}, []string{"grpc_type", "grpc_service", "grpc_method", "grpc_code"}),

		serverMsgSizeBytesSent: prom.NewCounterVec(
			prom.CounterOpts{
				Name: "grpc_server_bytes",
				Help: "Number of message sizes sent by the server.",
			}, []string{"grpc_type", "grpc_service", "grpc_method", "grpc_code"}),

		serverMsgSizeReceivedHistogramEnabled: false,
		serverMsgSizeReceivedHistogramOpts: prom.HistogramOpts{
			Name:    "grpc_server_msg_size_received_bytes",
			Help:    "Histogram of message sizes received by the server.",
			Buckets: defMsgBytesBuckets,
		},

		serverMsgSizeReceivedHistogram:    nil,
		serverMsgSizeSentHistogramEnabled: false,
		serverMsgSizeSentHistogramOpts: prom.HistogramOpts{
			Name:    "grpc_server_msg_size_sent_bytes",
			Help:    "Histogram of message sizes sent by the server.",
			Buckets: defMsgBytesBuckets,
		},
		serverMsgSizeSentHistogram: nil,
	}
}

// EnableMsgSizeReceivedBytesHistogram turns on recording of received message size of RPCs.
// Histogram metrics can be very expensive for Prometheus to retain and query. It takes
// options to configure histogram options such as the defined buckets.
func (m *ServerByteMetrics) EnableMsgSizeReceivedBytesHistogram() {
	m.serverMsgSizeReceivedHistogram = prom.NewHistogramVec(
		m.serverMsgSizeSentHistogramOpts,
		[]string{"grpc_service", "grpc_method", "grpc_stats"},
	)
	m.serverMsgSizeReceivedHistogramEnabled = true
}

// EnableMsgSizeSentBytesHistogram turns on recording of sent message size of RPCs.
// Histogram metrics can be very expensive for Prometheus to retain and query. It takes
// options to configure histogram options such as the defined buckets.
func (m *ServerByteMetrics) EnableMsgSizeSentBytesHistogram() {
	m.serverMsgSizeSentHistogram = prom.NewHistogramVec(
		m.serverMsgSizeSentHistogramOpts,
		[]string{"grpc_service", "grpc_method", "grpc_stats"},
	)
	m.serverMsgSizeSentHistogramEnabled = true
}

// Describe sends the super-set of all possible descriptors of metrics
// collected by this Collector to the provided channel and returns once
// the last descriptor has been sent.
func (m *ServerByteMetrics) Describe(ch chan<- *prom.Desc) {
	m.serverMsgSizeBytesReceived.Describe(ch)
	m.serverMsgSizeBytesSent.Describe(ch)
	if m.serverMsgSizeReceivedHistogramEnabled {
		m.serverMsgSizeReceivedHistogram.Describe(ch)
	}
	if m.serverMsgSizeSentHistogramEnabled {
		m.serverMsgSizeSentHistogram.Describe(ch)
	}
}

// Collect is called by the Prometheus registry when collecting
// metrics. The implementation sends each collected metric via the
// provided channel and returns once the last metric has been sent.
func (m *ServerByteMetrics) Collect(ch chan<- prom.Metric) {
	m.serverMsgSizeBytesReceived.Collect(ch)
	m.serverMsgSizeBytesSent.Collect(ch)
	if m.serverMsgSizeReceivedHistogramEnabled {
		m.serverMsgSizeReceivedHistogram.Collect(ch)
	}
	if m.serverMsgSizeSentHistogramEnabled {
		m.serverMsgSizeSentHistogram.Collect(ch)
	}
}

func (m *ServerByteMetrics) NewServerByteStatsHandler() stats.Handler {
	return &serverByteStatsHandler{
		serverByteMetrics: m,
	}
}

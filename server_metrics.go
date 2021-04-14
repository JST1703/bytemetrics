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
		serverMsgSizeBytesSent: prom.NewCounterVec(
			prom.CounterOpts{
				Name: "serverMsgSizeBytesSent",
				Help: "Total number of bytes sent by server.",
			}, []string{"grpc_service", "grpc_method", "grpc_stats"}),

		serverMsgSizeBytesReceived: prom.NewCounterVec(
			prom.CounterOpts{
				Name: "serverMsgSizeBytesReceived",
				Help: "Total number of bytes received by server.",
			}, []string{"grpc_service", "grpc_method", "grpc_stats"}),

		serverMsgSizeReceivedHistogramEnabled: false,
		serverMsgSizeReceivedHistogramOpts: prom.HistogramOpts{
			Name:    "grpc_server_msg_size_received_bytes",
			Help:    "Histogram of message sizes received by the server.",
			Buckets: defMsgBytesBuckets,
		},
		serverMsgSizeReceivedHistogram: nil,

		serverMsgSizeSentHistogramEnabled: false,
		serverMsgSizeSentHistogramOpts: prom.HistogramOpts{
			Name:    "grpc_server_msg_size_sent_bytes",
			Help:    "Histogram of message sizes sent by the server.",
			Buckets: defMsgBytesBuckets,
		},
		serverMsgSizeSentHistogram: nil,
	}
}

func (m *ServerByteMetrics) EnableMsgSizeReceivedBytesHistogram() {
	if !m.serverMsgSizeReceivedHistogramEnabled {
		m.serverMsgSizeReceivedHistogram = prom.NewHistogramVec(
			m.serverMsgSizeReceivedHistogramOpts,
			[]string{"grpc_service", "grpc_method", "grpc_stats"},
		)
	}
	m.serverMsgSizeReceivedHistogramEnabled = true
}

func (m *ServerByteMetrics) EnableMsgSizeSentBytesHistogram() {
	if !m.serverMsgSizeSentHistogramEnabled {
		m.serverMsgSizeSentHistogram = prom.NewHistogramVec(
			m.serverMsgSizeSentHistogramOpts,
			[]string{"grpc_service", "grpc_method", "grpc_stats"},
		)
	}
	m.serverMsgSizeSentHistogramEnabled = true
}

func (m *ServerByteMetrics) Describe(ch chan<- *prom.Desc) {
	m.serverMsgSizeBytesSent.Describe(ch)
	m.serverMsgSizeBytesReceived.Describe(ch)
	if m.serverMsgSizeReceivedHistogramEnabled {
		m.serverMsgSizeReceivedHistogram.Describe(ch)
	}
	if m.serverMsgSizeSentHistogramEnabled {
		m.serverMsgSizeSentHistogram.Describe(ch)
	}
}

func (m *ServerByteMetrics) Collect(ch chan<- prom.Metric) {
	m.serverMsgSizeBytesSent.Collect(ch)
	m.serverMsgSizeBytesReceived.Collect(ch)
	if m.serverMsgSizeReceivedHistogramEnabled {
		m.serverMsgSizeReceivedHistogram.Collect(ch)
	}
	if m.serverMsgSizeSentHistogramEnabled {
		m.serverMsgSizeSentHistogram.Collect(ch)
	}
}

func (m *ServerByteMetrics) NewServerByteStatsHandler() stats.Handler {
	return &ServerByteStatsHandler{
		serverByteMetrics: m,
	}
}

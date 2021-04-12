package bytesize

import (
	"time"
)

type serverByteReporter struct {
	metrics     *ServerByteMetrics
	rpcType     grpcType
	serviceName string
	methodName  string
	startTime   time.Time
}

func NewServerByteReporter(startTime time.Time, m *ServerByteMetrics, fullMethod string) *serverByteReporter {
	r := &serverByteReporter{
		metrics:   m,
		rpcType:   Unary,
		startTime: startTime,
	}
	r.serviceName, r.methodName = SplitMethodName(fullMethod)
	return r
}

// ReceivedMessageSize counts the size of received messages on server-side
func (r *serverByteReporter) ReceivedMessageSize(rpcStats grpcStats, size float64) {
	r.metrics.serverMsgSizeBytesReceived.WithLabelValues(r.serviceName, r.methodName, rpcStats.String()).Add(size)
	if r.metrics.serverMsgSizeReceivedHistogramEnabled {
		r.metrics.serverMsgSizeReceivedHistogram.WithLabelValues(r.serviceName, r.methodName, rpcStats.String()).Observe(size)
	}
}

// SentMessageSize counts the size of sent messages on server-side
func (r *serverByteReporter) SentMessageSize(rpcStats grpcStats, size float64) {
	r.metrics.serverMsgSizeBytesSent.WithLabelValues(r.serviceName, r.methodName, rpcStats.String()).Add(size)
	if r.metrics.serverMsgSizeSentHistogramEnabled {
		r.metrics.serverMsgSizeSentHistogram.WithLabelValues(r.serviceName, r.methodName, rpcStats.String()).Observe(size)
	}
}

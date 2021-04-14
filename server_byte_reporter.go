package bytesize

type ServerByteReporter struct {
	metrics     *ServerByteMetrics
	rpcType     grpcType
	serviceName string
	methodName  string
}

func NewServerByteReporter(m *ServerByteMetrics, fullMethod string) *ServerByteReporter {
	r := &ServerByteReporter{
		metrics: m,
		rpcType: Unary,
	}
	r.serviceName, r.methodName = SplitMethodName(fullMethod)
	return r
}

func (r *ServerByteReporter) ReceivedMessageSize(rpcStats grpcStats, size float64) {
	r.metrics.serverMsgSizeBytesReceived.WithLabelValues(r.serviceName, r.methodName, rpcStats.String()).Add(size)
	if r.metrics.serverMsgSizeReceivedHistogramEnabled {
		r.metrics.serverMsgSizeReceivedHistogram.WithLabelValues(r.serviceName, r.methodName, rpcStats.String()).Observe(size)
	}
}

func (r *ServerByteReporter) SentMessageSize(rpcStats grpcStats, size float64) {
	r.metrics.serverMsgSizeBytesSent.WithLabelValues(r.serviceName, r.methodName, rpcStats.String()).Add(size)
	if r.metrics.serverMsgSizeSentHistogramEnabled {
		r.metrics.serverMsgSizeSentHistogram.WithLabelValues(r.serviceName, r.methodName, rpcStats.String()).Observe(size)
	}
}

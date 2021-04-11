package bytesize

import (
	"context"

	"google.golang.org/grpc/stats"
)

type serverByteStatsHandler struct {
	serverByteMetrics *ServerByteMetrics
}

// TagRPC implements the stats.Hanlder interface.
func (h *serverByteStatsHandler) TagRPC(ctx context.Context, info *stats.RPCTagInfo) context.Context {
	rpcInfo := newRPCInfo(info.FullMethodName)
	return context.WithValue(ctx, &rpcInfoKey, rpcInfo)
}

// HandleRPC implements the stats.Hanlder interface.
func (h *serverByteStatsHandler) HandleRPC(ctx context.Context, s stats.RPCStats) {
	v, ok := ctx.Value(&rpcInfoKey).(*rpcInfo)
	if !ok {
		return
	}
	monitor := NewServerByteReporter(v.startTime, h.serverByteMetrics, v.fullMethodName)
	switch s := s.(type) {
	case *stats.InPayload:
		monitor.ReceivedMessageSize(Payload, float64(len(s.Data)))
	case *stats.OutPayload:
		monitor.SentMessageSize(Payload, float64(len(s.Data)))
	}
}

// TagConn implements the stats.Hanlder interface.
func (h *serverByteStatsHandler) TagConn(ctx context.Context, info *stats.ConnTagInfo) context.Context {
	return ctx
}

// HandleConn implements the stats.Hanlder interface.
func (h *serverByteStatsHandler) HandleConn(ctx context.Context, s stats.ConnStats) {
}

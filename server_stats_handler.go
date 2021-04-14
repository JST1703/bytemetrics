package bytesize

import (
	"context"

	"google.golang.org/grpc/stats"
)

type ServerByteStatsHandler struct {
	serverByteMetrics *ServerByteMetrics
}

// TagRPC implements the stats.Hanlder interface.
func (h *ServerByteStatsHandler) TagRPC(ctx context.Context, info *stats.RPCTagInfo) context.Context {
	rpcInfo := newRPCInfo(info.FullMethodName)
	return context.WithValue(ctx, &rpcInfoKey, rpcInfo)
}

// HandleRPC implements the stats.Hanlder interface.
func (h *ServerByteStatsHandler) HandleRPC(ctx context.Context, s stats.RPCStats) {
	v, ok := ctx.Value(&rpcInfoKey).(*rpcInfo)
	if !ok {
		return
	}
	monitor := NewServerByteReporter(h.serverByteMetrics, v.fullMethodName)
	switch s := s.(type) {
	case *stats.InPayload:
		monitor.ReceivedMessageSize(Payload, float64(len(s.Data)))
	case *stats.OutPayload:
		monitor.SentMessageSize(Payload, float64(len(s.Data)))
	}
}

// TagConn implements the stats.Hanlder interface.
func (h *ServerByteStatsHandler) TagConn(ctx context.Context, info *stats.ConnTagInfo) context.Context {
	return ctx
}

// HandleConn implements the stats.Hanlder interface.
func (h *ServerByteStatsHandler) HandleConn(ctx context.Context, s stats.ConnStats) {
}

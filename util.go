package bytesize

import (
	"strings"
	"time"
)

type grpcType string

const (
	Unary        grpcType = "unary"
	ClientStream grpcType = "client_stream"
	ServerStream grpcType = "server_stream"
	BidiStream   grpcType = "bidi_stream"
)

var (
	rpcInfoKey = "rpc-info"

	defMsgBytesBuckets = []float64{0, 32, 64, 128, 256, 512, 1024, 2048, 4096, 8192}
)

func SplitMethodName(fullMethodName string) (string, string) {
	fullMethodName = strings.TrimPrefix(fullMethodName, "/") // remove leading slash
	if i := strings.Index(fullMethodName, "/"); i >= 0 {
		return fullMethodName[:i], fullMethodName[i+1:]
	}
	return "unknown", "unknown"
}

type grpcStats string

const (
	// Header indicates that the stats is the header
	Header grpcStats = "header"

	// Payload indicates that the stats is the Payload
	Payload grpcStats = "payload"

	// Tailer indicates that the stats is the Payload
	Tailer grpcStats = "tailer"
)

// String function returns the grpcStats with string format.
func (s grpcStats) String() string {
	return string(s)
}

type rpcInfo struct {
	fullMethodName string
	startTime      time.Time
}

func newRPCInfo(fullMethodName string) *rpcInfo {
	return &rpcInfo{fullMethodName: fullMethodName}
}

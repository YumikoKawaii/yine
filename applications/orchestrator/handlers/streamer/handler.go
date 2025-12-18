package streamer

import (
	api "github.com/YumikoKawaii/rpc.com/protobuf/orchestrator"
	"github.com/YumikoKawaii/shared/logger"
	"google.golang.org/grpc"
)

type Handler struct {
	api.StreamerServer
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) ReceiveMessages(request *api.ReceiveMessagesRequest, stream grpc.ServerStreamingServer[api.Message]) error {
	logger.Infof("Stream opened for receiving messages")

	// TODO: Implement message streaming logic
	logger.Warnf("ReceiveMessages not yet implemented")

	return nil
}

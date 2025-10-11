package streamer

import (
	api "github.com/YumikoKawaii/rpc.com/protobuf/orchestrator"
	"google.golang.org/grpc"
)

type Handler struct {
	api.StreamerServer
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) ReceiveMessages(request *api.ReceiveMessagesRequest, stream grpc.ServerStreamingServer[api.Message]) error {
	return nil
}

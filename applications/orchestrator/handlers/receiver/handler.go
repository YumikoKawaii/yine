package receiver

import (
	"context"

	api "github.com/YumikoKawaii/rpc.com/protobuf/orchestrator"
)

type Handler struct {
	api.ReceiverServer
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) SendMessage(ctx context.Context, request *api.SendMessageRequest) (*api.SendMessageResponse, error) {
	return nil, nil
}

//func (h *Handler) ReceiveMessages(request *api.ReceiveMessagesRequest, stream grpc.ServerStreamingServer[api.Message]) error {
//	return nil
//}

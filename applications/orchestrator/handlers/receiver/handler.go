package receiver

import (
	"context"
	"net/http"
	"time"

	api "github.com/YumikoKawaii/rpc.com/protobuf/orchestrator"
	"github.com/YumikoKawaii/shared/pubsub"
	"github.com/golang/protobuf/proto"
	"github.com/samber/lo"
	"yumiko_kawaii.com/yine/applications/orchestrator/handlers/connection_registry"
	"yumiko_kawaii.com/yine/applications/orchestrator/pkg/constants"
	"yumiko_kawaii.com/yine/applications/orchestrator/pkg/models"
	"yumiko_kawaii.com/yine/applications/orchestrator/pkg/repository"
	"yumiko_kawaii.com/yine/applications/orchestrator/pkg/repository/uow"
)

type Handler struct {
	api.ReceiverServer
	connRegistry     connection_registry.Registry
	messagePublisher pubsub.Publisher
	worker           uow.IWorker
}

func NewHandler(connRegistry connection_registry.Registry, messagePublisher pubsub.Publisher, worker uow.IWorker) *Handler {
	return &Handler{
		connRegistry:     connRegistry,
		messagePublisher: messagePublisher,
		worker:           worker,
	}
}

func (h *Handler) SendMessage(ctx context.Context, request *api.SendMessageRequest) (*api.SendMessageResponse, error) {
	if err := h.worker.Do(ctx, func(store uow.IStore) error {
		if _, err := store.Messages().Upsert(ctx, &models.Message{
			Sender:         request.Sender,
			ConversationId: request.ConversationId,
			Content:        request.Content,
			Type:           request.Type.String(),
		}); err != nil {
			return err
		}

		userConversations, err := store.UserConversations().List(ctx, repository.UserConversationFilter{
			ConversationId: &request.ConversationId,
		})
		if err != nil {
			return err
		}

		userIdentifications := make([]string, constants.Zero)
		lo.ForEach(userConversations, func(item models.UserConversation, _ int) {
			userIdentifications = append(userIdentifications, item.UserIdentification)
		})
		servers, err := h.connRegistry.GetServers(userIdentifications)
		if err != nil {
			return err
		}
		// publish messages to servers
		// consider sync or async ?
		// just sync, this message is publish to the streamer, not the connection
		messageBytes, err := proto.Marshal(&api.Message{
			Sender:         request.Sender,
			ConversationId: request.ConversationId,
			Content:        request.Content,
			Type:           request.Type,
			Timestamp:      time.Now().Unix(),
			//Status:         0,
		})
		if err != nil {
			return err
		}
		for _, sv := range servers {
			if err := h.messagePublisher.Publish(ctx, constants.GenerateMessagesTopic(sv), messageBytes); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return &api.SendMessageResponse{
		Code:    int32(http.StatusOK),
		Message: "Success",
	}, nil
}

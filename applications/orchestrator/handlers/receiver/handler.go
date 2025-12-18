package receiver

import (
	"context"
	"net/http"
	"time"

	api "github.com/YumikoKawaii/rpc.com/protobuf/orchestrator"
	"github.com/YumikoKawaii/shared/logger"
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

func NewHandler(registry connection_registry.Registry, publisher pubsub.Publisher, worker uow.IWorker) *Handler {
	return &Handler{
		connRegistry:     registry,
		messagePublisher: publisher,
		worker:           worker,
	}
}

func (h *Handler) SendMessage(ctx context.Context, request *api.SendMessageRequest) (*api.SendMessageResponse, error) {
	logger.WithFields(logger.Fields{
		"sender":          request.Sender,
		"conversation_id": request.ConversationId,
		"message_type":    request.Type.String(),
	}).Infof("SendMessage request received")

	if err := h.worker.Do(ctx, func(store uow.IStore) error {
		if _, err := store.Messages().Upsert(ctx, &models.Message{
			Sender:         request.Sender,
			ConversationId: request.ConversationId,
			Content:        request.Content,
			Type:           request.Type.String(),
		}); err != nil {
			logger.WithFields(logger.Fields{
				"error":           err,
				"conversation_id": request.ConversationId,
			}).Errorf("Failed to upsert message")
			return err
		}

		userConversations, err := store.UserConversations().List(ctx, repository.UserConversationFilter{
			ConversationId: &request.ConversationId,
		})
		if err != nil {
			logger.WithFields(logger.Fields{
				"error":           err,
				"conversation_id": request.ConversationId,
			}).Errorf("Failed to list user conversations")
			return err
		}

		userIdentifications := make([]string, constants.Zero)
		lo.ForEach(userConversations, func(item models.UserConversation, _ int) {
			userIdentifications = append(userIdentifications, item.UserIdentification)
		})

		servers, err := h.connRegistry.GetServers(ctx, userIdentifications)
		if err != nil {
			logger.WithFields(logger.Fields{
				"error":           err,
				"conversation_id": request.ConversationId,
			}).Errorf("Failed to get connected servers")
			return err
		}

		messageBytes, err := proto.Marshal(&api.Message{
			Sender:         request.Sender,
			ConversationId: request.ConversationId,
			Content:        request.Content,
			Type:           request.Type,
			Timestamp:      time.Now().Unix(),
		})
		if err != nil {
			logger.WithFields(logger.Fields{
				"error":           err,
				"conversation_id": request.ConversationId,
			}).Errorf("Failed to marshal message")
			return err
		}

		for _, sv := range servers {
			topic := constants.GenerateMessagesTopic(sv)
			if err := h.messagePublisher.Publish(ctx, topic, messageBytes); err != nil {
				logger.WithFields(logger.Fields{
					"error":           err,
					"server":          sv,
					"conversation_id": request.ConversationId,
				}).Errorf("Failed to publish message")
				return err
			}
		}

		return nil
	}); err != nil {
		logger.WithFields(logger.Fields{
			"error":           err,
			"conversation_id": request.ConversationId,
		}).Errorf("SendMessage failed")
		return nil, err
	}

	logger.WithFields(logger.Fields{
		"conversation_id": request.ConversationId,
	}).Infof("SendMessage completed successfully")

	return &api.SendMessageResponse{
		Code:    int32(http.StatusOK),
		Message: "Success",
	}, nil
}

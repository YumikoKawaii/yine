package connection_registry

import (
	"context"

	"github.com/YumikoKawaii/shared/logger"
	"github.com/go-redis/redis/v8"
)

type Registry interface {
	Register(ctx context.Context, userIdentification string, serverIdentification string) error
	GetServers(ctx context.Context, userIdentifications []string) ([]string, error)
}

func NewRegistry(client *redis.Client) Registry {
	return &redisImpl{
		redisCli: client,
	}
}

type redisImpl struct {
	redisCli *redis.Client
}

func (i *redisImpl) Register(ctx context.Context, userIdentification string, serverIdentification string) error {
	return i.redisCli.SAdd(ctx, userIdentification, serverIdentification).Err()
}

func (i *redisImpl) GetServers(ctx context.Context, userIdentifications []string) ([]string, error) {
	servers := make([]string, 0)

	for _, id := range userIdentifications {
		svs, err := i.redisCli.SMembers(ctx, id).Result()
		if err != nil {
			logger.WithFields(logger.Fields{
				"error":               err,
				"user_identification": id,
			}).Errorf("Failed to get servers from Redis for user, continuing")
			continue
		}
		servers = append(servers, svs...)
	}

	return servers, nil
}

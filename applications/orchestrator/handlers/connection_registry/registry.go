package connection_registry

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	// userToServerPrefix is the Redis key prefix for user-to-server mapping
	userToServerPrefix = "user:server:"
	// defaultTTL is the default TTL for user-to-server mappings (24 hours)
	defaultTTL = 24 * time.Hour
)

type Registry interface {
	Register(userIdentification string, serverIdentification string) error
	GetServers(userIdentifications []string) ([]string, error)
	Unregister(userIdentification string) error
	Close() error
}

type redisImpl struct {
	client *redis.Client
	ttl    time.Duration
}

// NewRedisRegistry creates a new Redis-based connection registry
func NewRedisRegistry(client *redis.Client) Registry {
	return &redisImpl{
		client: client,
		ttl:    defaultTTL,
	}
}

// NewRedisRegistryWithTTL creates a new Redis-based connection registry with custom TTL
func NewRedisRegistryWithTTL(client *redis.Client, ttl time.Duration) Registry {
	return &redisImpl{
		client: client,
		ttl:    ttl,
	}
}

// Register maps a user identification to a server identification in Redis
func (r *redisImpl) Register(userIdentification string, serverIdentification string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	key := fmt.Sprintf("%s%s", userToServerPrefix, userIdentification)
	if err := r.client.Set(ctx, key, serverIdentification, r.ttl).Err(); err != nil {
		return fmt.Errorf("failed to register user %s to server %s: %w", userIdentification, serverIdentification, err)
	}
	return nil
}

// GetServers retrieves the list of unique server identifications for the given user identifications
func (r *redisImpl) GetServers(userIdentifications []string) ([]string, error) {
	if len(userIdentifications) == 0 {
		return []string{}, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Build keys for all users
	keys := make([]string, len(userIdentifications))
	for i, userID := range userIdentifications {
		keys[i] = fmt.Sprintf("%s%s", userToServerPrefix, userID)
	}

	// Use pipeline for efficient batch retrieval
	pipe := r.client.Pipeline()
	cmds := make([]*redis.StringCmd, len(keys))
	for i, key := range keys {
		cmds[i] = pipe.Get(ctx, key)
	}

	if _, err := pipe.Exec(ctx); err != nil && err != redis.Nil {
		return nil, fmt.Errorf("failed to get servers for users: %w", err)
	}

	// Collect unique server identifications
	serverSet := make(map[string]struct{})
	for _, cmd := range cmds {
		serverID, err := cmd.Result()
		if err == redis.Nil {
			// User not found in registry, skip
			continue
		}
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve server: %w", err)
		}
		serverSet[serverID] = struct{}{}
	}

	// Convert set to slice
	servers := make([]string, 0, len(serverSet))
	for serverID := range serverSet {
		servers = append(servers, serverID)
	}

	return servers, nil
}

// Unregister removes a user-to-server mapping from Redis
func (r *redisImpl) Unregister(userIdentification string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	key := fmt.Sprintf("%s%s", userToServerPrefix, userIdentification)
	if err := r.client.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("failed to unregister user %s: %w", userIdentification, err)
	}
	return nil
}

// Close closes the Redis client connection
func (r *redisImpl) Close() error {
	return r.client.Close()
}

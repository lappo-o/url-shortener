package mocks

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/mock"
)

type RedisClientMock struct {
	mock.Mock
}

func (m *RedisClientMock) Get(ctx context.Context, key string) *redis.StringCmd {
	args := m.Called(ctx, key)
	cmd := redis.NewStringCmd(ctx)
	if args.Get(0) == nil {
		cmd.SetErr(redis.Nil)
	} else {
		cmd.SetVal(args.String(0))
	}
	return cmd
}

func (m *RedisClientMock) Set(
	ctx context.Context,
	key string,
	value any,
	expiration time.Duration,
) *redis.StatusCmd {
	args := m.Called(ctx, key, value, expiration)
	cmd := redis.NewStatusCmd(ctx)
	if args.Get(0) != nil {
		cmd.SetErr(args.Error(0))
	} else {
		cmd.SetVal("OK")
	}
	return cmd
}

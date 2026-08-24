package idgen

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

const (
	CounterKey = "url:counter"
	Offset     = uint64(15_000_000)
)

type Generator struct {
	rdb *redis.Client
}

func Dial(ctx context.Context, redisURL string) (*Generator, error) {
	opt, err := redis.ParseURL(redisURL)

	if err != nil {
		return nil, fmt.Errorf("Redis URL: %w", err)
	}

	rdb := redis.NewClient(opt)

	if err := rdb.Ping(ctx).Err(); err != nil {
		_ = rdb.Close()

		return nil, fmt.Errorf("Redis ping: %w", err)
	}

	return New(rdb), nil
}

func New(rdb *redis.Client) *Generator {
	return &Generator{rdb: rdb}
}

func (g *Generator) Close() error {
	return g.rdb.Close()
}

func (g *Generator) Init(ctx context.Context) error {
	if err := g.rdb.SetNX(ctx, CounterKey, Offset, 0).Err(); err != nil {
		return fmt.Errorf("Init counter: %w", err)
	}

	return nil
}

func (g *Generator) Next(ctx context.Context) (uint64, error) {
	first, _, err := g.NextN(ctx, 1)

	if err != nil {
		return 0, err
	}

	return first, nil
}

func (g *Generator) NextN(ctx context.Context, n int) (uint64, uint64, error) {
	if n < 1 {
		return 0, 0, fmt.Errorf("Incr counter: n must be at least 1")
	}

	last, err := g.rdb.IncrBy(ctx, CounterKey, int64(n)).Result()

	if err != nil {
		return 0, 0, fmt.Errorf("Incr counter: %w", err)
	}

	first := uint64(last) - uint64(n) + 1

	return first, uint64(last), nil
}

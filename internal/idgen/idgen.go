package idgen

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

const (
	CounterKey = "url:counter"
	Offset     = uint64(62 * 62 * 62 * 62 * 62) // 62^5 = 916132832
)

type Generator struct {
	rdb *redis.Client
}

func New(rdb *redis.Client) *Generator {
	return &Generator{rdb: rdb}
}

func (g *Generator) Init(ctx context.Context) error {
	if err := g.rdb.SetNX(ctx, CounterKey, Offset, 0).Err(); err != nil {
		return fmt.Errorf("Init counter: %w", err)
	}

	return nil
}

func (g *Generator) Next(ctx context.Context) (uint64, error) {
	n, err := g.rdb.Incr(ctx, CounterKey).Result()

	if err != nil {
		return 0, fmt.Errorf("Incr counter: %w", err)
	}

	return uint64(n), nil
}

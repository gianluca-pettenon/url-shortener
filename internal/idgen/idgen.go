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

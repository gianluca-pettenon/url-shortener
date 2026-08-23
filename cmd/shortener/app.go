package main

import (
	"context"
	"os"

	"github.com/gianluca-pettenon/url-shortener/internal/db"
	"github.com/gianluca-pettenon/url-shortener/internal/idgen"
	"github.com/gianluca-pettenon/url-shortener/internal/urls"
)

func open(ctx context.Context) (*urls.Service, func(), error) {
	pool, err := db.NewPool(ctx, os.Getenv("DATABASE_URL"))

	if err != nil {
		return nil, nil, err
	}

	rdb, err := db.NewRedis(ctx, os.Getenv("REDIS_URL"))

	if err != nil {
		pool.Close()

		return nil, nil, err
	}

	ids := idgen.New(rdb)

	if err := ids.Init(ctx); err != nil {
		pool.Close()
		_ = rdb.Close()

		return nil, nil, err
	}

	return urls.NewService(ids, urls.NewRepository(pool), os.Getenv("HASHIDS_SALT")), func() {
		pool.Close()
		_ = rdb.Close()
	}, nil
}

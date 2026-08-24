package main

import (
	"context"
	"os"

	"github.com/gianluca-pettenon/url-shortener/internal/idgen"
	"github.com/gianluca-pettenon/url-shortener/internal/postgres"
	"github.com/gianluca-pettenon/url-shortener/internal/urls"
)

func open(ctx context.Context) (*urls.Service, func(), error) {
	pool, err := postgres.NewPool(ctx, os.Getenv("DATABASE_URL"))

	if err != nil {
		return nil, nil, err
	}

	ids, err := idgen.Dial(ctx, os.Getenv("REDIS_URL"))

	if err != nil {
		pool.Close()

		return nil, nil, err
	}

	if err := ids.Init(ctx); err != nil {
		pool.Close()
		_ = ids.Close()

		return nil, nil, err
	}

	svc, err := urls.NewService(ids, urls.NewRepository(pool), os.Getenv("HASHIDS_SALT"))

	if err != nil {
		pool.Close()
		_ = ids.Close()

		return nil, nil, err
	}

	return svc, func() {
		pool.Close()
		_ = ids.Close()
	}, nil
}

func openRead(ctx context.Context) (*urls.Service, func(), error) {
	pool, err := postgres.NewPool(ctx, os.Getenv("DATABASE_URL"))

	if err != nil {
		return nil, nil, err
	}

	svc, err := urls.NewService(nil, urls.NewRepository(pool), os.Getenv("HASHIDS_SALT"))

	if err != nil {
		pool.Close()

		return nil, nil, err
	}

	return svc, pool.Close, nil
}

func shortURL(code string) string {
	return os.Getenv("APP_DOMAIN") + "/" + code
}

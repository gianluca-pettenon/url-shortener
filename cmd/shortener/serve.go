package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gianluca-pettenon/url-shortener/internal/urls"
	"github.com/spf13/cobra"
)

type resolver interface {
	Resolve(ctx context.Context, code string) (string, error)
}

func serveCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "serve",
		Short: "Redirect short codes over HTTP",
		RunE: func(cmd *cobra.Command, args []string) error {
			svc, cleanup, err := openRead(cmd.Context())

			if err != nil {
				return err
			}

			defer cleanup()

			return listen(cmd.Context(), svc)
		},
	}
}

func listen(ctx context.Context, svc resolver) error {
	addr := ":" + os.Getenv("APP_PORT")

	mux := http.NewServeMux()
	mux.HandleFunc("GET /{code}", redirectHandler(svc))

	srv := &http.Server{Addr: addr, Handler: mux, ReadHeaderTimeout: 5 * time.Second}

	go func() {
		<-ctx.Done()
		shutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = srv.Shutdown(shutdown)
	}()

	log.Printf("Listening on %s", addr)

	err := srv.ListenAndServe()

	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}

	return err
}

func redirectHandler(svc resolver) http.HandlerFunc {
	var cache sync.Map

	return func(w http.ResponseWriter, r *http.Request) {
		code := r.PathValue("code")

		if v, ok := cache.Load(code); ok {
			redirect(w, r, v.(string))
			return
		}

		original, err := svc.Resolve(r.Context(), code)

		if errors.Is(err, urls.ErrInvalidCode) || errors.Is(err, urls.ErrNotFound) {
			http.NotFound(w, r)
			return
		}

		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		cache.Store(code, original)
		redirect(w, r, original)
	}
}

func redirect(w http.ResponseWriter, r *http.Request, url string) {
	w.Header().Set("Cache-Control", "public, max-age=31536000")
	http.Redirect(w, r, url, http.StatusFound)
}

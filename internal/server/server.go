package server

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"github.com/twirapp/language-processor/internal/detector"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In
	LC fx.Lifecycle

	Mux      *http.ServeMux
	Detector *detector.Detector
}

type handlers struct {
	detector *detector.Detector
}

func New(opts Opts) {
	s := &handlers{
		detector: opts.Detector,
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "3012"
	}

	opts.Mux.HandleFunc("GET /detect", s.detect)
	opts.Mux.HandleFunc("GET /detect/languages", s.languages)
	opts.Mux.HandleFunc("POST /translate", s.translate)

	server := &http.Server{
		Handler: opts.Mux,
		Addr:    "0.0.0.0:" + port,
	}

	opts.LC.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			slog.Info("Starting listening", slog.String("port", port))
			go func() {
				if err := server.ListenAndServe(); err != nil {
					panic(err)
				}
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			return server.Shutdown(ctx)
		},
	})
}

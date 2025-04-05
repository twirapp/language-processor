package main

import (
	"github.com/twirapp/language-processor/internal/detector"
	"github.com/twirapp/language-processor/internal/server"
	"go.uber.org/fx"
	"net/http"
)

func main() {
	fx.New(
		fx.NopLogger,
		fx.Provide(
			http.NewServeMux,
			detector.New,
		),
		fx.Invoke(server.New),
	).Run()
}

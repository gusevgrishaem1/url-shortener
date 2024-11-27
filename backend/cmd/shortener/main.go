package main

import (
	"context"
	"github.com/gusevgrishaem1/url-shortener/internal/shortener/server"
)

func main() {
	if err := RunApp(); err != nil {
		panic(err)
	}
}

func RunApp() error {
	ctx := context.Background()
	return server.StartServer(ctx)
}

package main

import (
	"context"
	server "github.com/gusevgrishaem1/url-shortener/internal/sever"
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

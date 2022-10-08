package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/ngrok/ngrok-go"
	"github.com/ngrok/ngrok-go/config"
)

func main() {
	if err := run(context.Background()); err != nil {
		log.Fatal(err)
	}
}

func run(ctx context.Context) error {
	tun, err := ngrok.StartTunnel(ctx,
		config.HTTPEndpoint(
			config.WithDomain("gophercon-github-demo.ngrok.io"),
			config.WithOAuth("github"),
		),
		ngrok.WithAuthtokenFromEnv(),
	)
	if err != nil {
		return err
	}
	fmt.Println("Listening at", tun.URL())
	return http.Serve(tun, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello gopher", r.Header.Get("ngrok-auth-user-name"))
	}))
}

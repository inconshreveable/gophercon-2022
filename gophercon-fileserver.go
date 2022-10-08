package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	ngrok "github.com/ngrok/ngrok-go"
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
			config.WithDomain("gophercon-fileserver-demo.ngrok.io"),
			config.WithOAuth("github"),
		),
		ngrok.WithAuthtokenFromEnv(),
	)
	if err != nil {
		return err
	}
	fmt.Println("Listening at", tun.URL())
	return http.Serve(tun, http.FileServer(http.Dir(".")))
}

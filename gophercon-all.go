// instructions
//
// mkdir ngrok-gophercon
// cd ngrok-gophercon
// go mod init example.com/gophercon
// go get -v
// go run gophercon-all.go
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/ngrok/ngrok-go"
	"github.com/ngrok/ngrok-go/config"
	"golang.org/x/sync/errgroup"
)

func main() {
	if err := run(context.Background()); err != nil {
		log.Fatal(err)
	}
}

func run(ctx context.Context) error {
	sess, err := ngrok.Connect(ctx, ngrok.WithAuthtokenFromEnv())
	if err != nil {
		return err
	}

	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		tun, err := sess.StartTunnel(ctx, config.HTTPEndpoint(
			config.WithDomain("gophercon-demo.ngrok.io"),
		))
		if err != nil {
			return err
		}
		fmt.Println("Listening at", tun.URL())
		return http.Serve(tun, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, "Hello Gophercon!")
		}))
	})
	g.Go(func() error {
		tun, err := sess.StartTunnel(ctx, config.HTTPEndpoint(
			config.WithDomain("gophercon-google-demo.ngrok.io"),
			config.WithOAuth("google"),
		))
		if err != nil {
			return err
		}
		fmt.Println("Listening at", tun.URL())
		return http.Serve(tun, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, "Hello gopher", r.Header.Get("ngrok-auth-user-name"))
		}))
	})
	g.Go(func() error {
		tun, err := sess.StartTunnel(ctx, config.HTTPEndpoint(
			config.WithDomain("gophercon-github-demo.ngrok.io"),
			config.WithOAuth("github"),
		))
		if err != nil {
			return err
		}
		fmt.Println("Listening at", tun.URL())
		return http.Serve(tun, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, "Hello gopher", r.Header.Get("ngrok-auth-user-name"))
		}))
	})
	g.Go(func() error {
		tun, err := sess.StartTunnel(ctx, config.HTTPEndpoint(
			config.WithDomain("gophercon-fileserver-demo.ngrok.io"),
			config.WithOAuth("github"),
		))
		if err != nil {
			return err
		}
		fmt.Println("Listening at", tun.URL())
		return http.Serve(tun, http.FileServer(http.Dir(".")))
	})

	return g.Wait()
}

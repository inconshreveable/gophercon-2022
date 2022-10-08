package main

import (
	"context"
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
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
	sess, err := ngrok.Connect(ctx, ngrok.WithAuthtokenFromEnv())
	if err != nil {
		return err
	}
	cas, err := loadCAs()
	if err != nil {
		return err
	}
	tun, err := sess.StartTunnel(ctx, config.HTTPEndpoint(
		config.WithDomain("gophercon-mtls-demo.ngrok.io"),
		config.WithMutualTLSCA(cas),
	))
	if err != nil {
		return err
	}
	return http.Serve(tun, http.HandlerFunc(handler))
}

func loadCAs() (*x509.Certificate, error) {
	rootPEM, err := ioutil.ReadFile("/home/ubuntu/ngrok/secrets/local/tls/root.crt.pem")
	if err != nil {
		return nil, err
	}
	rootDER, _ := pem.Decode(rootPEM)
	cert, err := x509.ParseCertificate(rootDER.Bytes)
	if err != nil {
		return nil, err
	}
	return cert, nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello block\n"))
}

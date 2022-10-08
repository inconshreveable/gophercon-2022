package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
)

func main() {
	if err := run(context.Background()); err != nil {
		log.Fatal(err)
	}
}

func run(ctx context.Context) error {
	l, err := net.Listen("tcp", ":8080")
	if err != nil {
		return err
	}
	fmt.Println("Listening at", l.Addr())
	return http.Serve(l, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello Gophercon!")
	}))
}

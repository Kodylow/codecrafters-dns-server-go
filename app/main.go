package main

import (
	"os"

	"github.com/codecrafters-io/dns-server-starter-go/app/server"
	"github.com/codecrafters-io/dns-server-starter-go/pkg/gotracer"
)

func main() {
	addr := "127.0.0.1:2053"

	log := gotracer.New()
	log.SetLevel(gotracer.LevelDebug)
	log.AddOutput(os.Stdout)

	srv := server.New(addr, log)

	if err := srv.Start(); err != nil {
		log.Error.Printf("Server error: %v", err)
	}
}

package main

import (
	"boards-merger/internal/web"
	"flag"
	"fmt"
)

func main() {
	port := flag.String("port", "8080", "Port number for the web server")
	flag.Parse()

	fmt.Printf("Starting web server on port %v", *port)
	if err := web.StartWebServer(*port); err != nil {
		fmt.Printf("Failed to start web server on port %v: %v", *port, err.Error())
	}
}

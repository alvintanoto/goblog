package main

import (
	"flag"

	"alvintanoto.id/blog/internal/route"
)

func main() {
	// Initialize Application Flags
	port := flag.String("port", ":3000", "Application Port")
	flag.Parse()

	// Start Server
	route.Init(*port)
}

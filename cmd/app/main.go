package main

import (
	"flag"

	"alvintanoto.id/blog/internal/database/connection"
	"alvintanoto.id/blog/internal/route"
)

func main() {
	// Initialize Application Flags
	port := flag.String("port", ":3000", "Application Port")
	dsn := flag.String("dsn", "", "Postgresql DSN")
	flag.Parse()

	new(connection.Postgresql).Init(*dsn)

	// Start Server
	route.Init(*port)
}

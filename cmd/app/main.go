package main

import (
	"flag"

	"alvintanoto.id/blog/internal/database/connection"
	"alvintanoto.id/blog/internal/route"
)

func main() {
	// Initialize Application Flags
	port := flag.String("port", ":3000", "Application Port")
	dsn := flag.String("dsn", "postgresql://alvintanoto:alvin2497@localhost:5432/blog", "Postgresql DSN")
	secret := flag.String("secret", "xT11TWwO60c*b3&*j42coY9eSPdzJ77W", "App Secret Key")
	flag.Parse()

	new(connection.Postgresql).Init(*dsn)

	// Start Server
	route.Init(*port, *secret)
}

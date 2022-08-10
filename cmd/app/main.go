package main

import "alvintanoto.id/blog/pkg/log"

func main() {
	logger := log.Get()
	logger.InfoLog.Println("App starting...")
}

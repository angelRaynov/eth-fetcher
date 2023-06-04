package main

import (
	"eth_fetcher/server"
	_ "github.com/lib/pq"
)

func main() {
	server.Run()
}

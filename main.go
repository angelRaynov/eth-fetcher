package main

import (
	"eth_fetcher/server"
	_ "github.com/lib/pq"
)



// TODO optimize table definitions
func main() {
	server.Run()


}



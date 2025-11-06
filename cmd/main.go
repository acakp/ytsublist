package main

import (
	"log"
	"ytsublist/internal/csv"
	"ytsublist/internal/web"
)

func main() {
	err := csv.AddChannel("channels.csv", "https://www.test.com")
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	web.Serve()
}

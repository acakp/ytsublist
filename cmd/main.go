package main

import (
	"log"
	"ytsublist/internal/csv"
	"ytsublist/internal/web"
)

func main() {
	err := csv.AddChannel("channels.csv", "https://youtube.com/@example")
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	web.Serve()
}

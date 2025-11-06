package web

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"ytsublist/internal/csv"
)

func Serve() {
	http.HandleFunc("/", handler)
	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	// reading csv
	chs, err := csv.ReadCsv("channels.csv")
	if err != nil {
		fmt.Println("error while reading csv:")
		log.Fatal(err)
	}
	// make youtube.com/@channel type links out of ids
	instance := "https://youtube.com"
	for i, ch := range chs {
		chs[i].ID = makeLink(instance, ch.ID)
	}
	// creating template
	tmpl, err := template.ParseFiles("templates/home.html")
	if err != nil {
		fmt.Println("error while creating html template:")
		log.Fatal(err)
	}
	tmpl.Execute(w, chs)
}

func makeLink(instance, id string) string {
	if id[0] == '@' {
		return instance + "/" + id
	} else {
		return instance + "/channel/" + id
	}
}

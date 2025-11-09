package web

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"ytsublist/internal/csv"
)

func Serve() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/add_channel", addChannelHandler)
	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", nil)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
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

func addChannelHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	parsedData := r.PostForm
	isAdded, err := csv.AddChannel("channels.csv", parsedData["chlink"][0])
	if !isAdded && err != nil {
		log.Println("can't find the channel id or the handle in the provided string")
	}
	if err != nil {
		log.Fatal(err)
	}
	r.PostForm.Add("isAdded", "true")
	http.Redirect(w, r, "/", 303)
}

func makeLink(instance, id string) string {
	if id[0] == '@' {
		return instance + "/" + id
	} else {
		return instance + "/channel/" + id
	}
}

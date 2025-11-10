package web

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"ytsublist/internal/csv"
)

type Flash struct {
	Message string
	Type    string // "success" or "fail"
}

type Data struct {
	Channels []csv.Channel
	Flash    *Flash
}

func Serve() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/add_channel", addChannelHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
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
	// checking if channel was added succesfully after addChannel
	status := r.URL.Query().Get("status")
	var flash Flash
	switch status {
	case "success":
		flash.Message = "Channel was added successfully"
		flash.Type = "success"
	case "fail":
		flash.Message = "Can't find channel url or handle in provided string"
		flash.Type = "fail"
	}
	d := Data{Channels: chs, Flash: &flash}
	tmpl.Execute(w, d)
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
	var status string
	if isAdded {
		status = "success"
	} else {
		status = "fail"
	}
	http.Redirect(w, r, "/?status="+status, 303)
}

func makeLink(instance, id string) string {
	if id[0] == '@' {
		return instance + "/" + id
	} else {
		return instance + "/channel/" + id
	}
}

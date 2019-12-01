package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const storagePath string = "./storage/"
const viewUrlPath string = "/view/"
const editUrlPath string = "/edit/"
const saveUrlPath string = "/save/"

type Page struct {
	Title string
	Body  []byte
}

func (p *Page) save() error {
	filename := storagePath + p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	filename := storagePath + title + ".txt"
	body, err := ioutil.ReadFile(filename)

	if err != nil {
		return nil, err
	}

	return &Page{
		Title: title,
		Body:  body,
	}, nil
}

func CheckIfError(err error) {
	if err != nil {
		_ = fmt.Errorf("error: %s", err.Error())
	}
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len(viewUrlPath):]
	p, _ := loadPage(title)
	_, _ = fmt.Fprintf(w, "<h1>%s</h1><div>%s</dir>", p.Title, p.Body)
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len(editUrlPath):]
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}

	_, _ = fmt.Fprintf(w, "<h1>Editing %s</h1>"+
		"<form action=\"/save/%s\" method=\"POST\">"+
		"<textarea name=\"body\">%s</textarea><br>"+
		"<input type=\"submit\" value=\"Save\">"+
		"</form>",
		p.Title, p.Title, p.Body)
}

func saveHandler(w http.ResponseWriter, r *http.Request) {

}

func main() {
	http.HandleFunc(viewUrlPath, viewHandler)
	http.HandleFunc(editUrlPath, editHandler)
	http.HandleFunc(saveUrlPath, saveHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

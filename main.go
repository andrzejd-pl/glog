package main

import (
	"fmt"
	"html/template"
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

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	t, err := template.ParseFiles(storagePath + tmpl + ".html")
	CheckIfError(err)

	_ = t.Execute(w, p)
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len(viewUrlPath):]
	p, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, editUrlPath+title, http.StatusFound)
		return
	}
	renderTemplate(w, "view", p)
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len(editUrlPath):]
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	renderTemplate(w, "edit", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len(saveUrlPath):]
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	p.save()
	http.Redirect(w, r, viewUrlPath+title, http.StatusFound)
}

func main() {
	http.HandleFunc(viewUrlPath, viewHandler)
	http.HandleFunc(editUrlPath, editHandler)
	http.HandleFunc(saveUrlPath, saveHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

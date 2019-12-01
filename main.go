package main

import (
	"errors"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
)

const storagePath string = "./storage/"
const viewUrlPath string = "/view/"
const editUrlPath string = "/edit/"
const saveUrlPath string = "/save/"
const validPathExpresion string = "^/(edit|save|view)/([a-zA-Z0-9]+)$"

type config struct {
	Templates *template.Template
	ValidPath *regexp.Regexp
}

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
		log.Println(err.Error())
	}
}

func renderTemplate(templates *template.Template, w http.ResponseWriter, tmpl string, p *Page) {
	CheckIfError(templates.ExecuteTemplate(w, tmpl+".html", p))
}

func getTitle(validPath *regexp.Regexp, w http.ResponseWriter, r *http.Request) (string, error) {
	m := validPath.FindStringSubmatch(r.URL.Path)
	if m == nil {
		http.NotFound(w, r)
		return "", errors.New("invalid page title")
	}

	return m[2], nil
}

func viewHandler(configuration *config) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		title, err := getTitle(configuration.ValidPath, writer, request)
		if err != nil {
			CheckIfError(err)
			return
		}

		p, err := loadPage(title)
		if err != nil {
			CheckIfError(err)
			http.Redirect(writer, request, editUrlPath+title, http.StatusFound)
			return
		}
		renderTemplate(configuration.Templates, writer, "view", p)
	}
}

func editHandler(configuration *config) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		title, err := getTitle(configuration.ValidPath, writer, request)
		if err != nil {
			CheckIfError(err)
			return
		}

		p, err := loadPage(title)
		if err != nil {
			CheckIfError(err)
			p = &Page{Title: title}
		}
		renderTemplate(configuration.Templates, writer, "edit", p)
	}
}

func saveHandler(configuration *config) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		title, err := getTitle(configuration.ValidPath, writer, request)
		if err != nil {
			CheckIfError(err)
			return
		}

		body := request.FormValue("body")
		p := &Page{Title: title, Body: []byte(body)}
		CheckIfError(p.save())

		http.Redirect(writer, request, viewUrlPath+title, http.StatusFound)
	}
}

func main() {
	configuration := &config{
		Templates: template.Must(template.ParseFiles(storagePath+"edit.html", storagePath+"view.html")),
		ValidPath: regexp.MustCompile(validPathExpresion),
	}
	http.HandleFunc(viewUrlPath, viewHandler(configuration))
	http.HandleFunc(editUrlPath, editHandler(configuration))
	http.HandleFunc(saveUrlPath, saveHandler(configuration))
	log.Fatal(http.ListenAndServe(":8080", nil))
}

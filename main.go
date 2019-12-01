package main

import (
	"fmt"
	"io/ioutil"
)

const storagePath string = "./storage/"

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

func main() {
	p1 := Page{Title: "TestPage", Body: []byte("This is a sample Page.")}
	err := p1.save()
	CheckIfError(err)

	p2, err := loadPage("TestPage")
	CheckIfError(err)

	fmt.Println(string(p2.Body))
}

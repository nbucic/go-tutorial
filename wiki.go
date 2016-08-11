package main

import (
	"io/ioutil"
	"fmt"
	"net/http"
	"strings"
	"html/template"
)

const (
	edit = "/edit/"
	save = "/save/"
	view = "/view/"
	templatesFolder = "templates/"
)

var templates = template.Must(template.ParseFiles(templatesFolder + "view.html", templatesFolder + "edit.html"))
type Page struct {
	Title string
	Body  []byte
}

func (p *Page) save() error {
	filename := p.Title + ".md"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	filename := title + ".md"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: strings.Replace(filename, ".md", "", -1), Body: body}, err
}

func throwError(w http.ResponseWriter, err string) {
	http.Error(w, err, http.StatusInternalServerError)
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(w, templatesFolder + tmpl + ".html", p)
	if err != nil {
		throwError(w, err.Error())
		return
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s", r.URL.Path[1:])
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/view/"):]
	p, err := loadPage(title)

	if err != nil {
		http.Redirect(w, r, edit + "/" + title, http.StatusFound)
		return
	}
	renderTemplate(w, "view", p)
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len(edit):]
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}

	renderTemplate(w, "edit", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len(save):]
	body := r.FormValue("body")

	p := &Page{Title: title, Body: []byte(body)}
	err := p.save()

	if err != nil {
		throwError(w, err.Error())
		return
	}

	http.Redirect(w, r, view + title, http.StatusFound)
}

func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/save/", saveHandler)
	http.HandleFunc("/edit/", editHandler)
	http.ListenAndServe(":8080", nil)
}

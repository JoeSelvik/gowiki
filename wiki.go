package main

import (
	"html/template"
	"io/ioutil"
	"net/http"
	"regexp"
	"errors"
)


type Page struct {
	Title string
	Body []byte
}

func (p *Page) save() error {
	filename := p.Title + ".txt"

	// WriteFile writes a slice of bytes to a text file
	// 0600 file created with read-write permissions for current user only
	return ioutil.WriteFile(filename, p.Body, 0600)
}



func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := ioutil.ReadFile(filename)

	if err != nil {
		return nil, err
	}

	return &Page{Title: title, Body: body}, nil
}



const lenPath = len("/view/")

// creates a template, panics if template given cannot be loaded
var templates = template.Must(template.ParseFiles("edit.html", "view.html"))

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}




func viewHandler(w http.ResponseWriter, r *http.Request) {
	title, err := getTitle(w, r)
	if err != nil {
		return
	}

	p, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, "edit/"+title, http.StatusFound)
		return
	}

	renderTemplate(w, "view", p)
}


func editHandler(w http.ResponseWriter, r *http.Request) {
	title, err := getTitle(w, r)
	if err != nil {
		return
	}

	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}

	renderTemplate(w, "edit", p)
}


func saveHandler(w http.ResponseWriter, r *http.Request) {
	title, err := getTitle(w, r)
	if err != nil {
		return
	}

	body := r.FormValue("body")

	// Must convert body, a string, into []byte
	p := &Page{Title: title, Body: []byte(body)}

	// Write the data to a file
	err = p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}



var titleValidator = regexp.MustCompile("^[a-zA-Z0-9]+$")

func getTitle(w http.ResponseWriter, r *http.Request) (title string, err error) {
	title = r.URL.Path[lenPath:]

	if !titleValidator.MatchString(title) {
		http.NotFound(w, r)
		err = errors.New("Invalid Page Title")
	}

	return
}



func main() {
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/edit/", editHandler)
	http.HandleFunc("/save/", saveHandler)
	http.ListenAndServe(":8000", nil)
}




package main

import (
	"encoding/json"
	"encoding/xml"
	"net/http"
	"path"
	"text/template"
)

func main() {
	http.HandleFunc("/", foo)
	http.HandleFunc("/text", renderPlainText)
	http.HandleFunc("/json", renderJson)
	http.HandleFunc("/xml", renderXml)
	http.HandleFunc("/file", serveFile)
	http.HandleFunc("/htmltemplate", renderHtmlTemplate)

	http.ListenAndServe(":3000", nil)
}

// send headers only
func foo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Server", "A go web server")
	w.WriteHeader(200)
}

// render plain text
func renderPlainText(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("RENDER_TEXT_OK"))
}

type Profile struct {
	Name    string
	Hobbies []string
}

// render JSON
func renderJson(w http.ResponseWriter, r *http.Request) {
	profile := Profile{"Alex", []string{"snowboarding", "programming"}}

	js, err := json.Marshal(profile)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

// render XML
func renderXml(w http.ResponseWriter, r *http.Request) {
	profile := Profile{"Alex", []string{"snowboarding", "programming"}}

	x, err := xml.MarshalIndent(profile, "", " ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/xml")
	w.Write(x)
}

// serve a file
func serveFile(w http.ResponseWriter, r *http.Request) {
	// Assuming you want to serve a photo at 'images/foo.txt'
	fp := path.Join("images", "foo.txt")
	http.ServeFile(w, r, fp)
}

// render HTML template
func renderHtmlTemplate(w http.ResponseWriter, r *http.Request) {
	profile := Profile{"Alex", []string{"snowboarding", "programming"}}

	fp := path.Join("templates", "index.html")
	tmpl, err := template.ParseFiles(fp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, profile); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

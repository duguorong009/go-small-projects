package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
)

type Person struct {
	Name string
	Age  int
}

func personCreate(w http.ResponseWriter, r *http.Request) {
	var p Person

	err := decodeJSONBody(w, r, &p)
	if err != nil {
		var mr *malformedRequest
		if errors.As(err, &mr) {
			http.Error(w, mr.msg, mr.status)
		} else {
			log.Print(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	fmt.Fprintf(w, "Person: %+v", p)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/person/create", personCreate)

	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}

package main

import (
	"errors"
	"io"
	"net/http"

	"example.com/fat-service-pattern/internal/service"
)

func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var input service.RegisteredUserInput

	err := app.decodeJSON(r.Body, &input)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	err = app.service.RegisterUser(&input)
	if err != nil {
		if errors.Is(err, service.ErrFailedValidation) {
			app.failedValidation(w, r, input.ValidationErrors)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (app *application) decodeJSON(body io.ReadCloser, input *service.RegisteredUserInput) error {

	return nil
}

func (app *application) badRequest(w http.ResponseWriter, r *http.Request, err error) {
	http.Error(w, err.Error(), http.StatusBadRequest)
}

func (app *application) failedValidation(w http.ResponseWriter, r *http.Request, err map[string]string) {
	errStr := ""
	for k, v := range err {
		errStr += k
		errStr += ":"
		errStr += v
		errStr += "\n"
	}
	http.Error(w, errStr, http.StatusBadRequest)
}

func (app *application) serverError(w http.ResponseWriter, r *http.Request, err error) {
	http.Error(w, err.Error(), http.StatusInternalServerError)
}

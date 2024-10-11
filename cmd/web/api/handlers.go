package api

import (
	"net/http"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("Home"))
	if err != nil {
		app.serverError(w, r, err)
	}
}

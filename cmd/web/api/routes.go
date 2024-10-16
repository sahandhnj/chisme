package api

import (
	"net/http"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("./static/"))

	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))
	mux.HandleFunc("GET /{$}", app.home)

	mux.HandleFunc("GET /mock/servers", getServers)
	mux.HandleFunc("GET /mock/applications", getApplications)
	mux.HandleFunc("GET /mock/resources", getResources)
	mux.HandleFunc("GET /mock/server/{server}", getServerByID)
	mux.HandleFunc("GET /mock/server/{server}/application/{application}", getApplicationByID)
	mux.HandleFunc("GET /mock/server/{server}/resource/{resource}", getResourceByID)

	return app.recoverPanic(app.logRequest(commonHeaders(mux)))
}

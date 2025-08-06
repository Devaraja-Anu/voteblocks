package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (app *application) routes() http.Handler {

	r := chi.NewRouter()
	r.Use(app.recoverPanic)

	r.Get("/healthcheck", app.healthCheck)

	r.Get("/v1/polls/{id}", app.getPollHandler)
	r.Get("/v1/polls", app.getAllPollsHandler)
	r.Post("/v1/polls", app.createPollhandler)

	return r
}

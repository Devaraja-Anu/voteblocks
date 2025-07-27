package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (app *application) routes() http.Handler {

	r := chi.NewRouter()
	r.Use(app.recoverPanic)

	r.Get("/healthcheck", app.healthCheck)
	// r.Get("/poll", app.getPolls)
	// r.Post("/poll", app.createPoll)
	// r.Post("/poll/{id}", app.postVote)

	return r
}

package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (app *application) routes() http.Handler {

	r := chi.NewRouter()
	r.Use(app.recoverPanic)

	r.Get("/healthcheck", app.healthCheck)

	r.Get("/v1/poll", app.getAllPolls)
	r.Get("/v1/poll/{id}", app.getPoll)
	r.Post("/v1/polls", app.createPoll)
	r.Post("/v1poll/{id}/vote", app.postVote)
	r.Get("/v1/poll/{id}/results", app.getPollResults)

	return r
}

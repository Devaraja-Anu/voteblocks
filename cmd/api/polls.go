package main

import (
	"context"
	"net/http"
	"time"

	"Github.com/Devaraja-Anu/voteblocks/internal/db"
	"Github.com/Devaraja-Anu/voteblocks/internal/validator"
	"github.com/jackc/pgx/v5/pgtype"
)

func (app *application) createPollhandler(w http.ResponseWriter, r *http.Request) {

	var input struct {
		Title   string     `json:"title"`
		Desc    string     `json:"description"`
		Options []string   `json:"options"`
		Expiry  *time.Time `json:"expiry"`
	}

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	poll := &db.CreatePollParams{
		Title:       input.Title,
		Description: input.Desc,
		Options:     input.Options,
	}

	if input.Expiry != nil {
		poll.ExpiresAt = pgtype.Timestamp{
			Time:  *input.Expiry,
			Valid: true,
		}
	} else {
		poll.ExpiresAt = pgtype.Timestamp{
			Time:  time.Now().Add(24 * time.Hour),
			Valid: true,
		} // NULL in DB
	}

	v := validator.New()

	if validator.ValidatePolls(v, poll); !v.Valid() {
		app.failedValidation(w, r, v.Errors)
		return
	}

	createdPoll, err := app.queries.CreatePoll(ctx, *poll)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"polls": createdPoll}, nil)
	if err != nil {
		app.serverError(w, r, err)
	}
}

func (app *application) getPollHandler(w http.ResponseWriter, r *http.Request) {

	id, err := app.readParamId(r)
	if err != nil {
		app.notFound(w, r)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	returnval, err := app.queries.GetPoll(ctx, id)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"poll": returnval}, nil)
	if err != nil {
		app.serverError(w, r, err)
	}
}

func (app *application) getAllPollsHandler(w http.ResponseWriter, r *http.Request) {

	// err := app.writeJSON(w, http.StatusCreated, envelope{"id": id}, nil)
	// if err != nil {
	// 	app.serverError(w, r, err)
	// }

	// ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	// defer cancel()

}

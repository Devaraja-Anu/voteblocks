package main

import (
	"context"
	"math"
	"net/http"
	"strconv"
	"time"

	"Github.com/Devaraja-Anu/voteblocks/internal/db"
	"Github.com/Devaraja-Anu/voteblocks/internal/types"
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
		poll.ExpiresAt = pgtype.Timestamptz{
			Time:  *input.Expiry,
			Valid: true,
		}
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

	err = app.writeJSON(w, http.StatusOK, envelope{"poll": returnval}, nil)
	if err != nil {
		app.serverError(w, r, err)
	}
}

func (app *application) getAllPollsHandler(w http.ResponseWriter, r *http.Request) {

	query := r.URL.Query()

	search := query.Get("search")
	page, err := strconv.Atoi(query.Get("page"))
	if err != nil || page < 1 {
		page = 1
	}
	pageSize, err := strconv.Atoi(query.Get("page_size"))
	if err != nil || pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize

	v := validator.New()
	v.Check(page > 0, "page", "must be greater than zero")
	v.Check(pageSize > 0, "page_size", "must be greater than zero")
	v.Check(pageSize <= 50, "page_size", "must be a maximum of 50")

	if !v.Valid() {
		app.failedValidation(w, r, v.Errors)
		return
	}

	params := db.ListPollsParams{
		PlaintoTsquery: search,
		Limit:          int32(pageSize),
		Offset:         int32(offset),
	}

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	pollsList, err := app.queries.ListPolls(ctx, params)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	var totalRecords int64
	if len(pollsList) > 0 {
		totalRecords = pollsList[0].TotalRecords
	}

	metadata := types.Metadata{
		CurrentPage:  page,
		PageSize:     pageSize,
		FirstPage:    1,
		LastPage:     int(math.Ceil(float64(totalRecords) / float64(pageSize))),
		TotalRecords: int(totalRecords),
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"metadata":metadata,"polls": pollsList}, nil)
	if err != nil {
		app.serverError(w, r, err)
	}
}

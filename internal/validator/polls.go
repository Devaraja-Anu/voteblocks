package validator

import (
	"time"

	"Github.com/Devaraja-Anu/voteblocks/internal/db"
)

func ValidatePolls(v *Validator, polls *db.CreatePollParams) {
	v.Check(polls.Title != "", "title", "Title must be provided")
	v.Check(len(polls.Options) >= 2, "options", "there must be at least 2 options")
	if polls.ExpiresAt.Valid {
		v.Check(polls.ExpiresAt.Time.After(time.Now().Add(time.Hour)),
			"expiry", "Expiry must be at least an hour from now")
	}
}

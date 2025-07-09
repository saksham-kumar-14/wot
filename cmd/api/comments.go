package main

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/saksham-kumar-14/wot/internal/store"
)

type CommentPayload struct {
	Content string `json:"content" validate:"required,max=10000"`
	PostId  int    `json:"post_id"`
}

func (app *application) postComment(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "postID")
	postID, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	var payload CommentPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	userId := 1
	comment := &store.Comment{
		Content: payload.Content,
		UserId:  int64(userId),
		PostId:  int64(postID),
	}
	ctx := r.Context()

	if err := app.store.Comments.CreateComment(ctx, comment); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := writeJSON(w, http.StatusCreated, comment); err != nil {
		app.internalServerError(w, r, err)
		return
	}

}

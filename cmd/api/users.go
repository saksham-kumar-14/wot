package main

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/saksham-kumar-14/wot/internal/store"
)

type userKey string

const userCtx userKey = "user"

func (app *application) userContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idParam := chi.URLParam(r, "userID")
		userID, err := strconv.ParseInt(idParam, 10, 64)
		if err != nil {
			app.internalServerError(w, r, err)
			return
		}
		ctx := r.Context()

		user, err := app.store.Users.GetByID(ctx, int(userID))
		if err != nil {
			switch {
			case errors.Is(err, store.ErrNotFound):
				app.notFoundError(w, r, err)
			default:
				app.internalServerError(w, r, err)
			}
			return
		}

		//
		ctx = context.WithValue(ctx, userCtx, user)
		next.ServeHTTP(w, r.WithContext(ctx))

	})
}

func getUserFromContext(r *http.Request) *store.User {
	user, _ := r.Context().Value(userCtx).(*store.User)
	return user
}

func (app *application) getUser(w http.ResponseWriter, r *http.Request) {
	user := getUserFromContext(r)

	if err := writeJSON(w, http.StatusOK, user); err != nil {
		app.internalServerError(w, r, err)
		return
	}

}

type FriendUser struct {
	UserId int `json:"user_id"`
}

func (app *application) friendHandler(w http.ResponseWriter, r *http.Request) {
	friendUser := getUserFromContext(r)

	// TODO: revert when auth impl
	var payload FriendUser
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	ctx := r.Context()
	if err := app.store.Friends.Friend(ctx, int(friendUser.ID), payload.UserId); err != nil {
		switch err {
		case store.ErrAlreadyExists:
			app.conflict(w, r, err)
		default:
			app.internalServerError(w, r, err)
			return
		}
	}

	if err := writeJSON(w, http.StatusOK, payload); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) unfriendHandler(w http.ResponseWriter, r *http.Request) {
	unfriendUser := getUserFromContext(r)

	// TODO: revert when auth impl
	var payload FriendUser
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	ctx := r.Context()
	if err := app.store.Friends.Unfriend(ctx, int(unfriendUser.ID), payload.UserId); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := writeJSON(w, http.StatusOK, payload); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

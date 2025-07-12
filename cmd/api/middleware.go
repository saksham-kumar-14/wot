package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/saksham-kumar-14/wot/internal/store"
)

func (app *application) AuthTokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			writeJSONError(w, http.StatusBadRequest, "Authorization header is missing")
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			writeJSONError(w, http.StatusBadRequest, "Authorization header if invalid")
			return
		}

		token := parts[1]
		jwtToken, err := app.authenticator.ValidateToken(token)
		if err != nil {
			app.badRequestError(w, r, err)
			return
		}

		claims, _ := jwtToken.Claims.(jwt.MapClaims)
		userID, err := strconv.ParseInt(fmt.Sprintf("%.f", claims["sub"]), 10, 64)
		if err != nil {
			app.badRequestError(w, r, errors.New("unauthorized"))
			return
		}

		ctx := r.Context()
		// user, err := app.store.Users.GetByID(ctx, int(userID))
		// if err != nil {
		// 	app.badRequestError(w, r, errors.New("unauthorized"))
		// 	return
		// }
		user, err := app.getUser(ctx, userID)
		if err != nil {
			app.badRequestError(w, r, err)
			return
		}

		ctx = context.WithValue(ctx, userCtx, user)
		next.ServeHTTP(w, r.WithContext(ctx))

	})
}

func (app *application) getUser(ctx context.Context, userID int64) (*store.User, error) {
	app.logger.Infow("cache hit", "key", "user", "id", userID)

	user, err := app.cacheStorage.Users.Get(ctx, userID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		app.logger.Infow("fetching from DB", "id", userID)
		user, err = app.store.Users.GetByID(ctx, int(userID))
		if err != nil {
			return nil, err
		}

		if err := app.cacheStorage.Users.Set(ctx, user); err != nil {
			return nil, err
		}
	}

	return user, nil
}

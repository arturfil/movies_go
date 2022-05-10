package main

import (
	"context"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func (app *application) wrap(next http.Handler) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		ctx := context.WithValue(r.Context(), "params", ps)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func (app *application) routes() http.Handler {
	router := httprouter.New()
	secure := alice.New(app.checkToken)
	router.HandlerFunc(http.MethodGet, "/api/status", app.statusHandler)
	router.HandlerFunc(http.MethodPost, "/api/v1/graphql/list", app.moviesGraphQL)
	router.HandlerFunc(http.MethodPost, "/api/v1/signin", app.SignIn)
	router.GET("/api/v1/movies", app.wrap(secure.ThenFunc(app.getAllMovies)))
	router.HandlerFunc(http.MethodGet, "/api/v1/movie/:id", app.getOneMovie)
	router.HandlerFunc(http.MethodGet, "/api/v1/genres", app.getAllGenres)
	router.HandlerFunc(http.MethodGet, "/api/v1/movies/genre/:genre_id", app.getMoviesByGenre)
	router.HandlerFunc(http.MethodPost, "/api/v1/admin/createMovie", app.createMovie)
	router.POST("/api/v1/admin/updateMovie", app.wrap(secure.ThenFunc(app.updatetMovie)))
	router.HandlerFunc(http.MethodDelete, "/api/v1/admin/deleteMovie/:id", app.deleteMovie)

	return app.enableCORS(router)
}

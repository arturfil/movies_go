package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()
	router.HandlerFunc(http.MethodGet, "/api/status", app.statusHandler)
	router.HandlerFunc(http.MethodPost, "/api/v1/signin", app.SignIn)
	router.HandlerFunc(http.MethodGet, "/api/v1/movies", app.getAllMovies)
	router.HandlerFunc(http.MethodGet, "/api/v1/movie/:id", app.getOneMovie)
	router.HandlerFunc(http.MethodGet, "/api/v1/genres", app.getAllGenres)
	router.HandlerFunc(http.MethodGet, "/api/v1/movies/genre/:genre_id", app.getMoviesByGenre)
	router.HandlerFunc(http.MethodPost, "/api/v1/admin/createMovie", app.createMovie)
	router.HandlerFunc(http.MethodPost, "/api/v1/admin/updateMovie", app.updatetMovie)
	router.HandlerFunc(http.MethodDelete, "/api/v1/admin/deleteMovie/:id", app.deleteMovie)

	return app.enableCORS(router)
}

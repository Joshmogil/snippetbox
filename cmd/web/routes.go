package main

import (
	"net/http"
	"github.com/justinas/alice"
	"github.com/bmizerany/pat"
)


func (app *application) routes() http.Handler {

	standardMiddleware := alice.New(app.recoverPanic, app.logRequest,secureHeaders)

	/*mux := http.NewServeMux()
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet", app.showSnippet)
	mux.HandleFunc("/snippet/create", app.createSnippet)
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	*/
	mux := pat.New()
	mux.Get("/", http.HandlerFunc(app.home))
	mux.Get("/snippet/create",http.HandlerFunc(app.createSnippetForm))
	mux.Post("/snippet/create", http.HandlerFunc(app.createSnippet))
	mux.Get("/snippet/:id", http.HandlerFunc(app.showSnippet)) // Moved down
	
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Get("/static/", http.StripPrefix("/static", fileServer))

	return standardMiddleware.Then(mux)
}
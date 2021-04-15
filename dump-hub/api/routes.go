package api

/*
The MIT License (MIT)
Copyright (c) 2021 Davide Pataracchia
Permission is hereby granted, free of charge, to any person
obtaining a copy of this software and associated documentation
files (the "Software"), to deal in the Software without
restriction, including without limitation the rights to use,
copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the
Software is furnished to do so, subject to the following
conditions:
The above copyright notice and this permission notice shall be
included in all copies or substantial portions of the Software.
THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
OTHER DEALINGS IN THE SOFTWARE.
*/

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

/*
defineRoutes :: Define API routes and handlers
*/
func (engine *Engine) defineRoutes() {
	router := mux.NewRouter().StrictSlash(true)
	router.Use(accessLogger)

	router.
		Name("Upload").
		Path(engine.baseAPI + "upload").
		Methods(http.MethodPost).
		HandlerFunc(upload(engine.eClient))

	router.
		Name("History").
		Path(engine.baseAPI + "history").
		Methods(http.MethodPost).
		HandlerFunc(getHistory(engine.eClient))

	router.
		Name("Search").
		Path(engine.baseAPI + "search").
		Methods(http.MethodPost).
		HandlerFunc(search(engine.eClient))

	router.
		Name("Delete").
		Path(engine.baseAPI + "delete").
		Methods(http.MethodPost).
		HandlerFunc(delete(engine.eClient))

	engine.router = router
}

/*
accessLogger :: Access log middleware
*/
func accessLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf(
			"(%s) %s %s",
			r.Host,
			r.UserAgent(),
			r.URL,
		)
		next.ServeHTTP(w, r)
	})
}

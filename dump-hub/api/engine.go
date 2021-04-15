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
	"strconv"

	"github.com/gorilla/mux"
	"github.com/x0e1f/dump-hub/elastic"
)

/*
Engine :: Core API Engine
*/
type Engine struct {
	host    string
	port    int
	baseAPI string
	router  *mux.Router
	eClient *elastic.Client
}

const pageSize = 20

/*
New :: Create the Api Engine object
*/
func New(host string, port int, baseAPI string, eClient *elastic.Client) *Engine {
	log.Println("Initializing engine...")
	engine := &Engine{
		host:    host,
		port:    port,
		baseAPI: baseAPI,
		eClient: eClient,
	}
	engine.defineRoutes()

	return engine
}

/*
Serve :: Serve API
*/
func (engine *Engine) Serve() {
	log.Printf("Serving API on %s:%d", engine.host, engine.port)

	addr := engine.host + ":" + strconv.Itoa(engine.port)
	log.Fatal(http.ListenAndServe(addr, engine.router))
}

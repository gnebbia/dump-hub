package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/x0e1f/dump-hub/elastic"
)

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

type historyReq struct {
	Page int `json:"page"`
}

/*
getHistory :: Get history documents
*/
func getHistory(eClient *elastic.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var historyReq historyReq

		err := json.NewDecoder(r.Body).Decode(&historyReq)
		if err != nil {
			log.Println(err)
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		from := pageSize * (historyReq.Page - 1)
		historyData, err := eClient.GetHistory(from, pageSize)
		if err != nil {
			log.Printf("(ERROR) (%s) %s", r.URL, err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
		response, err := json.Marshal(historyData)
		if err != nil {
			http.Error(w, "", http.StatusInternalServerError)
			log.Println(err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(response)
	}
}

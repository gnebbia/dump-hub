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
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"log"
	"net/http"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/x0e1f/dump-hub/common"
	"github.com/x0e1f/dump-hub/elastic"
	"github.com/x0e1f/dump-hub/parser"
)

const chunkSize = 1000

/*
upload :: Upload dump file (POST)
*/
func upload(eClient *elastic.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseMultipartForm(1024 * 1024)

		/* Get Pattern Value */
		pattern := r.FormValue("pattern")
		if len(pattern) <= 0 {
			log.Printf("(ERROR) (%s) pattern value not found", r.URL)
			http.Error(w, "", http.StatusBadRequest)
			return
		}
		/* Get Columns Value */
		columns := r.FormValue("columns")
		if len(columns) <= 0 {
			log.Printf("(ERROR) (%s) columns value not found", r.URL)
			http.Error(w, "", http.StatusBadRequest)
		}

		p, err := parser.New(pattern, columns)
		if err != nil {
			log.Printf("(ERROR) Parser creation error: (%s)", err)
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		/* Get File Body */
		file, handler, err := r.FormFile("file")
		if err != nil {
			log.Printf("(ERROR) (%s) %s", r.URL, err)
			http.Error(w, "", http.StatusBadRequest)
			return
		}
		defer file.Close()

		/* Check Content-Type */
		contentType := handler.Header["Content-Type"][0]
		if !strings.Contains(contentType, "text/") {
			log.Printf("(ERROR) (%s) invalid content type.", r.URL)
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		/* Upload file on tmp */
		filePath := "/tmp/" + uuid.New().String()
		tmpFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			log.Println(err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
		defer tmpFile.Close()
		io.Copy(tmpFile, file)

		/* Compute file checksum */
		f, err := os.Open(filePath)
		if err != nil {
			log.Println(err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
		hash := sha256.New()
		if _, err := io.Copy(hash, f); err != nil {
			log.Println(err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
		checkSum := hex.EncodeToString(hash.Sum(nil))

		/* Check if file already exist */
		fileExist, err := eClient.IsAlreadyUploaded(checkSum)
		if err != nil {
			log.Println(err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
		if fileExist {
			log.Println("File already exist: ", checkSum)
			http.Error(w, "", http.StatusTeapot)
			return
		}

		/* Create history document */
		date := time.Now().Format("2006-01-02 15:04:05")
		history := common.History{
			Date:     date,
			Filename: handler.Filename,
			Checksum: checkSum,
			Status:   0,
		}
		err = eClient.NewHistory(&history, checkSum)
		if err != nil {
			log.Println(err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		/* Process file in background */
		go processFile(
			eClient,
			p,
			handler.Filename,
			filePath,
			checkSum,
		)

		w.WriteHeader(http.StatusOK)
	}
}

/*
processFile :: Process file line by line
*/
func processFile(e *elastic.Client, p *parser.Parser, fn string, fp string, cs string) {
	/* Open file from tmp */
	file, err := os.Open(fp)
	if err != nil {
		e.UpdateHistoryStatus(cs, -1)
		log.Println(err)
	}
	defer func() {
		file.Close()
		err := os.Remove(fp)
		if err != nil {
			log.Println(err)
		}
	}()

	/* Start uploader routines */
	var wg sync.WaitGroup
	quitChan := make(chan struct{})
	entryChan := make(chan map[string]string)
	for i := 0; i < runtime.GOMAXPROCS(0); i++ {
		go uploader(i, &wg, e, quitChan, entryChan)
	}

	/* Scan file line by line */
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		/* Parse entry document */
		entry := p.ParseEntry(fn, cs, scanner.Text())
		if entry == nil {
			continue
		}
		entryChan <- entry
	}

	close(quitChan)
	close(entryChan)

	wg.Wait()
	log.Printf("Processing complete: %s", fn)

	/* Refresh elastic index */
	e.Refresh()
	/* Update history status (Complete)*/
	e.UpdateHistoryStatus(cs, 1)
}

/*
uploader :: Upload entries to elastic
*/
func uploader(id int, wg *sync.WaitGroup, e *elastic.Client, quitChan <-chan struct{}, entryChan <-chan map[string]string) {
	wg.Add(1)
	log.Printf("Starting uploader #%d\n", id)
	run := true
	chunk := []map[string]string{}

	for run {
		/* Chunk size reached */
		if len(chunk) >= chunkSize {
			err := e.BulkInsert(chunk)
			if err != nil {
				log.Println(err)
			}
			chunk = []map[string]string{}
		}

		select {
		case <-quitChan:
			run = false
		case entry := <-entryChan:
			if entry == nil {
				continue
			}
			/* Append document to chunk */
			chunk = append(chunk, entry)
		}
	}

	/* If there is still data, upload chunk */
	if len(chunk) > 0 {
		err := e.BulkInsert(chunk)
		if err != nil {
			log.Println(err)
		}
	}

	wg.Done()
	log.Printf("Shutting down uploader #%d\n", id)
}

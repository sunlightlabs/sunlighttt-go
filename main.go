/* {{{ Copyright (c) Paul R. Tagliamonte <paultag@gmail.com>, 2015
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in
 * all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
 * THE SOFTWARE. }}} */

package main

import (
	"fmt"
	"os"

	"encoding/json"
	"net/http"
)

func writeJSON(w http.ResponseWriter, data interface{}, code int) error {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		return err
	}
	return nil
}

func writeError(w http.ResponseWriter, message string, code int) error {
	return writeJSON(w, map[string][]map[string]string{
		"errors": []map[string]string{
			map[string]string{"message": message},
		},
	}, code)
}

func isValidKey(key string) bool {
	validKey := os.Getenv("IFTTTChannelKey")
	if validKey == "" {
		return true
	}

	if key == validKey {
		return true
	}

	return false
}

func main() {
	server := NewTriggerServer()
	server.Register("/ifttt/v1/triggers/upcoming-bills", UpcomingBillTrigger{})
	server.Register("/ifttt/v1/triggers/new-bills-query", NewBillsQueryTrigger{})
	server.Register("/ifttt/v1/triggers/new-laws", NewLawsTrigger{})
	// "/ifttt/v1/triggers/new-legislators"
	// "/ifttt/v1/triggers/congress-birthdays"

	mux := http.NewServeMux()

	mux.HandleFunc("/ifttt/v1/status", func(w http.ResponseWriter, req *http.Request) {
		iftttChannelKey := req.Header.Get("IFTTT-Channel-Key")
		if !isValidKey(iftttChannelKey) {
			writeError(w, "Bad key given", 401)
			return
		}
		fmt.Fprintf(w, "Hello, World")
	})

	mux.HandleFunc("/ifttt/v1/test/setup", func(w http.ResponseWriter, req *http.Request) {
		iftttChannelKey := req.Header.Get("IFTTT-Channel-Key")
		if !isValidKey(iftttChannelKey) {
			writeError(w, "Bad key given", 401)
			return
		}
		fmt.Fprintf(w, `{
        "data": { "samples": { "triggers": {
			"new-bills-query": {"query": "\"Common Core\""},
			"new-legislators": {
				"location": {
					"lat": "44.967586",
					"lon": "-103.772234",
					"address": "19424 Us Highway 85, Belle Fourche, SD 57717",
					"description": "Geographic Center of the United States"
				}}}}}}`)
	})

	mux.HandleFunc("/ifttt/v1/triggers/", func(w http.ResponseWriter, req *http.Request) {
		iftttChannelKey := req.Header.Get("IFTTT-Channel-Key")
		if !isValidKey(iftttChannelKey) {
			writeError(w, "Bad key given", 401)
			return
		}

		fmt.Printf("Handling request for %s\n", req.URL.Path)
		response, err := server.Handle(req)
		if err != nil {
			writeError(w, err.Error(), 400)
			return
		}
		if response == nil {
			writeError(w, "Internal error", 500)
			return
		}
		writeJSON(w, response, 200)
	})
	http.ListenAndServe(":8000", mux)
}

// vim: foldmethod=marker

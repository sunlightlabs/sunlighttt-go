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
	"encoding/json"
	"errors"
	"net/http"
)

type TriggerFields struct {
	TriggerFields struct {
		Query    string `json:"query"`
		Location struct {
			Lat         string `json:"lat"`
			Lon         string `json:"lon"`
			Address     string `json:"address"`
			Description string `json:"description"`
		}
	} `json:"triggerFields"`

	User  map[string]string `json:"user"`
	Limit int               `json:"limit"`
}

type Trigger interface {
	Handle(TriggerFields) (interface{}, error)
}

type TriggerServer struct {
	handlers map[string]Trigger
}

func NewTriggerServer() TriggerServer {
	return TriggerServer{
		handlers: map[string]Trigger{},
	}
}

type Meta struct {
	Id        string `json:"id"`
	Timestamp int64  `json:"timestamp"`
}

func (trigger *TriggerServer) Register(path string, hook Trigger) {
	trigger.handlers[path] = hook
}

func (trigger *TriggerServer) Handle(req *http.Request) (interface{}, error) {
	/* OK, first of all, let's get the POST body */
	if req.Method != "POST" {
		return nil, errors.New("We only support POST at this time.")
	}

	/* Now, let's parse the POST body, which is a JSON data payload */
	data := TriggerFields{}
	data.Limit = -1
	err := json.NewDecoder(req.Body).Decode(&data)
	if err != nil {
		return nil, err
	}

	/* We're good. Let's look up the route and fetch our handler. */
	if val, ok := trigger.handlers[req.URL.Path]; ok {
		/* Dispatch to the registered hook. */
		return val.Handle(data)
	}
	return nil, errors.New("No such trigger")
}

// vim: foldmethod=marker

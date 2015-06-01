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
	"errors"
	"sort"

	"github.com/sunlightlabs/go-sunlight/congress"
)

type NewLegislatorsData struct {
	Date            string `json:"date"`
	Name            string `json:"name"`
	Party           string `json:"party"`
	Phone           string `json:"phone"`
	State           string `json:"state"`
	TwitterUsername string `json:"twitter_username"`
	Website         string `json:"website"`
	Meta            Meta   `json:"meta"`
}

type NewLegislators struct {
	Data []NewLegislatorsData `json:"data"`
}

func (b NewLegislators) Len() int {
	return len(b.Data)
}

func (b NewLegislators) Less(i, j int) bool {
	return b.Data[i].Meta.Timestamp > b.Data[j].Meta.Timestamp
	// return b.Data[i].birthday.Before(b.Data[j].birthday)
}

func (b NewLegislators) Swap(i, j int) {
	b.Data[i], b.Data[j] = b.Data[j], b.Data[i]
}

type NewLegislatorsTrigger struct {
	Trigger
}

func (trigger NewLegislatorsTrigger) Handle(fields TriggerFields) (interface{}, error) {
	location := fields.TriggerFields.Location
	if location == nil {
		return nil, errors.New("Need a lat/lon")
	}

	ret := NewLegislators{
		Data: make([]NewLegislatorsData, 0),
	}

	people, err := congress.GetLegislatorsByLatLon(
		location.Lat,
		location.Lon,
		map[string]string{
			"fields": "title,first_name,last_name,bioguide_id,state,party,district,term_start,twitter_id,phone,website",
		},
	)

	if err != nil {
		return nil, err
	}

	for _, person := range people.Results {
		start, err := parseTime(person.TermStart)
		if err != nil {
			return nil, err
		}

		entry := NewLegislatorsData{
			Name:            personName(&person),
			Party:           person.Party,
			Phone:           person.Phone,
			State:           person.State,
			TwitterUsername: person.TwitterId,
			Website:         person.Website,
			Date:            person.TermStart,
			Meta: Meta{
				Id:        person.BioguideId + "/" + person.State,
				Timestamp: start.Unix(),
			},
		}
		ret.Data = append(ret.Data, entry)
	}

	sort.Sort(ret)

	if fields.Limit != -1 {
		ret.Data = ret.Data[:fields.Limit]
	}

	return ret, nil
}

// vim: foldmethod=marker

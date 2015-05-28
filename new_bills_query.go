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

	"github.com/sunlightlabs/go-sunlight/congress"
)

type NewBillsQueryData struct {
	Code            string `json:"code"`
	IntroducedOn    string `json:"introduced_on"`
	OfficialURL     string `json:"official_url"`
	OpenCongressURL string `json:"open_congress_url"`
	Query           string `json:"query"`
	SponsorName     string `json:"sponsor_name"`
	Title           string `json:"title"`
	Meta            Meta   `json:"meta"`
	Date            string `json:"date"`
}

type NewBillsQuery struct {
	Data []NewBillsQueryData `json:"data"`
}

type NewBillsQueryTrigger struct {
	Trigger
}

func (trigger NewBillsQueryTrigger) Handle(fields TriggerFields) (interface{}, error) {
	query := fields.TriggerFields.Query
	if query == "" {
		return nil, errors.New("Missing a trigger field `query`")
	}

	bills, err := congress.BillTextSearch(
		query,
		map[string]string{
			"fields": "bill_id,bill_type,number,introduced_on,short_title,official_title,sponsor,urls.congress,urls.opencongress",
			"order":  "congress,introduced_on,number",
		},
	)
	if err != nil {
		return nil, err
	}

	ret := NewBillsQuery{}
	ret.Data = make([]NewBillsQueryData, 0)

	for _, bill := range bills.Results {
		if fields.Limit != -1 && len(ret.Data) >= fields.Limit {
			break
		}

		dtime, err := parseTime(bill.IntroducedOn)
		if err != nil {
			return nil, err
		}

		ret.Data = append(ret.Data, NewBillsQueryData{
			Code:            displayBillId(&bill),
			Query:           query,
			Date:            bill.IntroducedOn,
			Title:           displayTitle(&bill),
			IntroducedOn:    displayDate(dtime),
			SponsorName:     sponsorName(&bill),
			OfficialURL:     bill.Urls["congress"],
			OpenCongressURL: bill.Urls["opencongress"],
			Meta: Meta{
				Id:        bill.BillId,
				Timestamp: dtime.Unix(),
			},
		})
	}

	return ret, nil
}

// vim: foldmethod=marker

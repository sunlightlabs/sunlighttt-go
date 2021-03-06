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
	"github.com/sunlightlabs/go-sunlight/congress"
)

type NewLawsData struct {
	BecameLawOn     string
	Code            string
	OfficialURL     string
	OpenCongressURL string
	SponsorName     string
	Title           string
	Meta            Meta   `json:"meta"`
	Date            string `json:"date"`
}

type NewLaws struct {
	Data []NewLawsData `json:"data"`
}

type NewLawsTrigger struct {
	Trigger
}

func (trigger NewLawsTrigger) Handle(fields TriggerFields) (interface{}, error) {
	bills, err := congress.BillSearch(
		map[string]string{
			"fields":          "bill_id,bill_type,number,history.enacted_at,short_title,official_title,sponsor,urls.congress,urls.opencongress",
			"order":           "history.enacted_at",
			"history.enacted": "true",
		},
	)
	if err != nil {
		return nil, err
	}

	ret := NewLaws{}
	ret.Data = make([]NewLawsData, 0)

	for _, bill := range bills.Results {
		if fields.Limit != -1 && len(ret.Data) >= fields.Limit {
			break
		}

		dtime, err := parseTime(bill.History.EnactedAt)
		if err != nil {
			return nil, err
		}

		ret.Data = append(ret.Data, NewLawsData{
			Code:            displayBillId(&bill),
			Date:            bill.History.EnactedAt,
			Title:           displayTitle(&bill),
			BecameLawOn:     displayDate(dtime),
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

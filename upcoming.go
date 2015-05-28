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

type UpcomingBillsData struct {
	Chamber         string
	Code            string
	LegislativeDate string
	SourceURL       string
	SponsorName     string
	Title           string
	Meta            Meta   `json:"meta"`
	Date            string `json:"date"`
}

type UpcomingBills struct {
	Data []UpcomingBillsData `json:"data"`
}

type UpcomingBillTrigger struct {
	Trigger
}

func (trigger UpcomingBillTrigger) Handle(fields TriggerFields) (interface{}, error) {
	bills, err := congress.UpcomingBills(map[string]string{
		"order":         "scheduled_at",
		"range__exists": "true",
		"fields":        "bill_id,chamber,legislative_day,range,url,bill,scheduled_at",
	})

	if err != nil {
		return nil, err
	}

	ret := UpcomingBills{}
	ret.Data = make([]UpcomingBillsData, 0)

	for _, bill := range bills.Results {
		if fields.Limit != -1 && len(ret.Data) >= fields.Limit {
			break
		}

		legislativeDay, err := parseTime(bill.LegislativeDay)
		if err != nil {
			continue
		}

		ret.Data = append(ret.Data, UpcomingBillsData{
			Code:            displayBillId(&bill.Bill),
			Title:           displayTitle(&bill.Bill),
			SponsorName:     sponsorName(&bill.Bill),
			SourceURL:       bill.URL,
			LegislativeDate: displayDate(legislativeDay),
			Date:            bill.LegislativeDay,
			Chamber:         displayChamber(bill.Chamber),
			Meta: Meta{
				Id:        bill.Range + "/" + bill.LegislativeDay + "/" + bill.BillId,
				Timestamp: legislativeDay.Unix(),
			},
		})
	}

	return ret, nil
}

// vim: foldmethod=marker

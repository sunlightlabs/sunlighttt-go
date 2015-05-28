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
	"strconv"
	"time"

	"github.com/sunlightlabs/go-sunlight/congress"
)

func displayTitle(bill *congress.Bill) string {
	if bill.ShortTitle != "" {
		return bill.ShortTitle
	}
	return bill.OfficialTitle
}

func displayChamber(chamber string) string {
	switch chamber {
	case "senate":
		return "Senate"
	case "house":
		return "House of Representatives"
	default:
		return "Unknown"
	}
}

func displayBillId(bill *congress.Bill) string {
	// Bill types can be: hr, hres, hjres, hconres, s, sres, sjres, sconres.

	var billType = "unknown"
	if val, ok := map[string]string{
		"hr":      "H.R.",
		"hres":    "H.Res.",
		"hjres":   "H.J.Res.",
		"hconres": "H.Con.Res.",
		"s":       "S.",
		"sres":    "S.R.",
		"sjres":   "S.J.Res.",
		"sconres": "S.Con.Res.",
	}[bill.BillType]; ok {
		billType = val
	}

	return billType + " " + strconv.Itoa(bill.Number)
}

func sponsorName(bill *congress.Bill) string {
	return bill.Sponsor.Title + ". " +
		bill.Sponsor.FirstName + " " + bill.Sponsor.LastName
}

func personName(who *congress.Legislator) string {
	return who.Title + ". " +
		who.FirstName + " " + who.LastName
}

func parseTime(date string) (time.Time, error) {
	est, err := time.LoadLocation("America/New_York")
	if err != nil {
		return time.Time{}, err
	}
	dateFormat := "2006-01-02"
	return time.ParseInLocation(dateFormat, date, est)
}

func displayDate(when time.Time) string {
	var suffix = "th"
	switch when.Day() % 10 {
	case 1:
		suffix = "st"
	case 2:
		suffix = "nd"
	case 3:
		suffix = "rd"
	}
	return when.Format("Jan 2" + suffix + ", 2006")
}

// vim: foldmethod=marker

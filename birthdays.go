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
	"sort"

	// "errors"
	"github.com/sunlightlabs/go-sunlight/congress"
	"time"
)

type BirthdaysData struct {
	birthday        time.Time
	Age             int    `json:"age"`
	BirthYear       int    `json:"birth_year"`
	BirthDay        string `json:"birthday"`
	BirthDayDate    string `json:"birthday_date"`
	Date            string `json:"date"`
	Name            string `json:"name"`
	Party           string `json:"party"`
	State           string `json:"state"`
	TwitterUsername string `json:"twitter_username"`
	Meta            Meta   `json:"meta"`
}

type Birthdays struct {
	Data []BirthdaysData `json:"data"`
}

func (b Birthdays) Len() int {
	return len(b.Data)
}

func (b Birthdays) Less(i, j int) bool {
	return b.Data[i].Meta.Timestamp < b.Data[j].Meta.Timestamp
	// return b.Data[i].birthday.Before(b.Data[j].birthday)
}

func (b Birthdays) Swap(i, j int) {
	b.Data[i], b.Data[j] = b.Data[j], b.Data[i]
}

type BirthdaysTrigger struct {
	Trigger
}

func (trigger BirthdaysTrigger) Handle(fields TriggerFields) (interface{}, error) {
	people, err := congress.GetLegislators(map[string]string{
		"fields":    "title,first_name,chamber,last_name,state,party,district,birthday,bioguide_id,twitter_id",
		"in_office": "true",
		"per_page":  "all",
	})

	if err != nil {
		return nil, err
	}

	ret := Birthdays{}
	ret.Data = make([]BirthdaysData, 0)
	today := time.Now()
	yesterday := today.Add(-(24 * time.Hour))

	for _, person := range people.Results {
		birthday, err := parseTime(person.Birthday)
		if err != nil {
			return nil, err
		}

		age := today.Year() - birthday.Year()

		/*
		 * The following monumental and horseshit hack due to the fact that
		 * Durations don't have a year method. Since leap years are hard
		 * and I don't want to implement that in a userland app, I'm going
		 * to jump through insane hoops to get it.
		 *
		 * Thanks for this, Go.
		 */
		cakeDay := time.Date(
			today.Year(),
			birthday.Month(),
			birthday.Day(),
			birthday.Hour(),
			0, 0, 0,
			birthday.Location(),
		)

		/* So, now we can figure out if they have already had their birthday */
		if cakeDay.After(today) {
			/* So, they've not had their birthday yet */
			age = age - 1
		} else if cakeDay.Before(yesterday) {
			/* Party's over; sorry! */
			continue
		}

		state := person.State
		if person.Chamber == "house" {
			state = fmt.Sprintf("%s-%d", person.State, person.District)
		}

		ret.Data = append(ret.Data, BirthdaysData{
			birthday:        cakeDay,
			Age:             age,
			BirthYear:       birthday.Year(),
			BirthDay:        displayDate(birthday),
			BirthDayDate:    person.Birthday,
			Date:            cakeDay.Format("2006-01-02"),
			Name:            personName(&person),
			Party:           person.Party,
			State:           state,
			TwitterUsername: person.TwitterId,
			Meta: Meta{
				Id:        fmt.Sprintf("%d/%s", today.Year(), person.BioguideId),
				Timestamp: cakeDay.Unix(),
			},
		})
	}

	sort.Sort(ret)

	if fields.Limit != -1 {
		ret.Data = ret.Data[:fields.Limit]
	}

	return ret, nil
}

// vim: foldmethod=marker

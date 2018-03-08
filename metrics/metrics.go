// Copyright 2018 National Library of Norway
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package metrics provides metrics for Sigridr
package metrics

import (
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	r "gopkg.in/gorethink/gorethink.v3"

	"github.com/nlnwa/sigridr/database"
)

type Sigridr struct {
	pattern string
	h       http.Handler
	Error   *log.Logger
	Db      *database.Rethink
}

// New allocates and returns a new metric for Sigridr
func New(db *database.Rethink, logger *log.Logger, pattern string) *Sigridr {
	reg := prometheus.NewRegistry()
	s := &Sigridr{
		pattern: pattern,
		h: promhttp.HandlerFor(reg, promhttp.HandlerOpts{
			ErrorLog: logger,
		}),
		Error: logger,
		Db: db,
	}

	c := prometheus.NewCounterFunc(prometheus.CounterOpts{
		Name: "sigridr_twitter_statuses_total",
		Help: "Total number of twitter statuses harvested",
	}, s.Total)

	reg.MustRegister(c)

	return s
}

func (s *Sigridr) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.Handle(s.pattern, s.h)
	return mux
}

// Total returns the total number of tweets harvested
//
// If database communication fails -1 is returned as a sentinel and the error is logged
func (s *Sigridr) Total() float64 {
	if err := s.Db.Connect(); err != nil {
		s.Error.Printf("%v", err)
		return -1
	}
	defer s.Db.Disconnect()

	cursor, err := r.Table("execution").Filter(func(e r.Term) r.Term {
		return e.HasFields("statuses")
	}).Field("statuses").Sum().Run(s.Db.Session)
	if err != nil {
		s.Error.Printf("%v", err)
		return -1
	}
	count := new(float64)
	if err := cursor.One(count); err != nil {
		s.Error.Printf("%v", err)
		return -1
	} else {
		return *count
	}
}

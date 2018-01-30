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

	r "gopkg.in/gorethink/gorethink.v3"

	"github.com/nlnwa/sigridr/database"
)

type Sigridr struct {
	Error *log.Logger
	Db    *database.Rethink
}

// New allocates and returns a new metric for Sigridr
func New(dbHost string, dbPort int, dbName string, logger *log.Logger) *Sigridr {
	return &Sigridr{
		Error: logger,
		Db:    database.New(database.WithName(dbName), database.WithAddress(dbHost, dbPort)),
	}
}

// Total returns the total number of tweets harvested
//
// If database communication fails -1 is returned as a sentinel and the error is logged
func (s *Sigridr) Total() interface{} {
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
	count := new(int)
	if err := cursor.One(count); err != nil {
		s.Error.Printf("%v", err)
		return -1
	} else {
		return *count
	}
}

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

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/namsral/flag"

	"github.com/nlnwa/sigridr-metrics/expvar"
	"github.com/nlnwa/sigridr-metrics/metrics"
)

func main() {
	port := flag.Int("port", 8081, "port to listen on")
	help := flag.Bool("help", false, "print this help message")
	dbPort := flag.Int("db-port", 28015, "database port")
	dbHost := flag.String("db-host", "localhost", "database host")
	dbName := flag.String("db-name", "sigridr", "database name")
	pattern := flag.String("pattern", "/metrics", "path of metrics endpoint")

	flag.Parse()

	if *help {
		flag.PrintDefaults()
		os.Exit(1)
	}

	logger := log.New(os.Stderr, "ERROR: ", log.LstdFlags)

	m := metrics.New(*dbHost, *dbPort, *dbName, logger)
	expvar.Publish("count", expvar.Func(m.Total))

	mux := http.NewServeMux()
	mux.Handle(*pattern, expvar.Handler())

	if err := http.ListenAndServe(fmt.Sprintf(":%d", *port), mux); err != nil {
		logger.Fatal(err)
	}
}
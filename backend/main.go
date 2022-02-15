// Copyright (c) 2020 Computing Infrastructure Research Centre (CIRC), McMaster
// University. All rights reserved.

package main

import (
	"time"

	"github.com/mcmaster-circ/canids-v2/backend/api"
	"github.com/mcmaster-circ/canids-v2/backend/auth"
	"github.com/mcmaster-circ/canids-v2/backend/libraries/scheduler"
	"github.com/mcmaster-circ/canids-v2/backend/state"
	log "github.com/sirupsen/logrus"
)

// gitHash is commit hash populated by Docker
var gitHash string

func main() {
	// initialize main state
	s, err := state.Provision(gitHash)
	if err != nil {
		s.Log.Fatal(err)
	}

	// enable debug logging if requested in configuration
	if s.Config.DebugLogging {
		s.Log.Info("[main] debug logging enabled")
		s.Log.SetLevel(log.DebugLevel)
	}

	// begin scheduled refreshing of alarm ip sets
	scheduler.Provision(map[string]string{
		"firehol_abusers_1d":  "https://iplists.firehol.org/files/firehol_abusers_1d.netset",
		"firehol_abusers_30d": "https://iplists.firehol.org/files/firehol_abusers_30d.netset",
		"firehol_anonymous":   "https://iplists.firehol.org/files/firehol_anonymous.netset",
		"firehol_level1":      "https://iplists.firehol.org/files/firehol_level1.netset",
		"firehol_level2":      "https://iplists.firehol.org/files/firehol_level2.netset",
		"firehol_level3":      "https://iplists.firehol.org/files/firehol_level3.netset",
	}, 18*time.Hour, s.AlarmManager)

	// provision API state
	a, err := auth.Provision(s)
	if err != nil {
		s.Log.Fatal(err)
	}

	// start API server
	err = api.Start(s, a)
	if err != nil {
		s.Log.Fatal(err)
	}
}

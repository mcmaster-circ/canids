// Copyright (c) 2020 Computing Infrastructure Research Centre (CIRC), McMaster
// University. All rights reserved.

package main

import (
	"os"
	"strconv"
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

	var skipScheduler bool
	skipScheduler = false
	skipSchedulerStr := os.Getenv("SKIP_SCHEDULER")
	skipScheduler, err := strconv.ParseBool(skipSchedulerStr)
	if err != nil {
		skipScheduler = false
	}

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
	if !skipScheduler {
		err = scheduler.Provision(s, 18*time.Hour, s.AlarmManager)
		if err != nil {
			s.Log.Fatal(err)
		}
	}

	// provision API state
	a, err := auth.Provision(s)
	if err != nil {
		s.Log.Fatal(err)
	}

	p := auth.ProvisionAuthPage(s)

	// start API server
	err = api.Start(s, a, p)
	if err != nil {
		s.Log.Fatal(err)
	}
}

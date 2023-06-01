// Copyright (c) 2020 Computing Infrastructure Research Centre (CIRC), McMaster
// University. All rights reserved.

// Package dashboard provides the dashboard API service for the backend.
package dashboard

import (
	"github.com/mcmaster-circ/canids-v2/backend/libraries/elasticsearch"
	"github.com/mcmaster-circ/canids-v2/backend/libraries/uuid"
	"github.com/mcmaster-circ/canids-v2/backend/state"
)

var (
	// defaultViews are the prepoulated views if no dashboard exists
	defaultViews = []elasticsearch.DocumentView{
		{
			Name:       elasticsearch.DefaultViewName,
			Class:      "bar",
			Fields:     []string{"timestamp"},
			FieldNames: []string{"Time"},
		},
	}
)

// provisionDashboard is called by the getHandler. It is called when a dashboard
// does not exist. Provision will query for the views. If there does not exist views,
// it will create new views based on the defaults above.
func provisionDashboard(s *state.State) (elasticsearch.DocumentDashboard, error) {
	var dash elasticsearch.DocumentDashboard

	// attempt to query existing views
	views, err := elasticsearch.AllView(s)
	if err != nil || len(views) == 0 {
		// load and save default views to database
		views = defaultViews
		for i := range views {
			// populate missing fields
			views[i].UUID = uuid.Generate()
			// save view in database
			_, err = views[i].Index(s)
			if err != nil {
				// failed to save view
				return dash, err
			}
		}
	}
	// either using existing views or new views, get list of view UUIDs
	viewUUIDs := []string{}
	viewSizes := []elasticsearch.SizeClass{}
	for _, view := range views {
		viewUUIDs = append(viewUUIDs, view.UUID)
		viewSizes = append(viewSizes, elasticsearch.SizeHalf)
	}
	// generate dashboard
	dash.UUID = uuid.Generate()
	dash.Name = "Main Dashboard"
	dash.Views = viewUUIDs
	dash.Sizes = viewSizes

	// save dashboard
	_, err = dash.Index(s)
	if err != nil {
		return dash, err
	}
	// successfully generated new dashboard
	return dash, nil
}

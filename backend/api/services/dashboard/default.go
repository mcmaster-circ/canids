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
	// defaultViews are the prepoulated views if a group does not have one
	defaultViews = []elasticsearch.DocumentView{
		{
			UUID:       "", // UUID populated below
			Group:      "", // Group populated below
			Authorized: "", // Authorized populated below
			Name:       elasticsearch.DefaultViewName,
			Class:      "bar",
			Fields:     []string{"timestamp"},
			FieldNames: []string{"Time"},
		},
	}
)

// provisionDashboard is called by the getHandler. It is called when a dashboard
// does not exist for a group. Provision will query for the group's views. If
// there are groups, the dashboard will be created using them. If there does not
// exist views for the group, it will create new views based on the defaults
// above. The authorized asset used will be the first asset registered in the
// group.
func provisionDashboard(s *state.State, group elasticsearch.DocumentGroup) (elasticsearch.DocumentDashboard, error) {
	var dash elasticsearch.DocumentDashboard

	// attempt to query existing views
	views, err := elasticsearch.QueryViewByGroup(s, group.UUID)
	if err != nil || len(views) == 0 {
		// load and save default views to database
		views = defaultViews
		// default authorized asset will be the first available one
		asset := group.Authorized[0]
		for i := range views {
			// populate missing fields
			views[i].UUID = uuid.Generate()
			views[i].Group = group.UUID
			views[i].Authorized = asset
			// add asset name to visualization
			views[i].Name += " " + asset
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
	dash.Group = group.UUID
	dash.Name = group.Name + " Dashboard"
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

package mPagerDuty

import "github.com/PagerDuty/go-pagerduty"

func getFakedOnCalls() *pagerduty.ListOnCallsResponse {
	var response = &pagerduty.ListOnCallsResponse{
		APIListObject: pagerduty.APIListObject{Limit: 100, Offset: 0, More: false, Total: 0},
		OnCalls: []pagerduty.OnCall{
			{
				User:             getFakedUser(),
				Schedule:         getFakedSchedule(),
				EscalationPolicy: getFakedEscalationPolicy(),
				EscalationLevel:  1,
				Start:            "2021-07-21T14:54:39Z",
				End:              "2022-12-30T14:04:11Z",
			},
		},
	}

	return response
}

func getFakedUser() pagerduty.User {
	user := pagerduty.User{
		APIObject: pagerduty.APIObject{
			ID:      "PJ6XOVE",
			Self:    "https://api.pagerduty.com/users/PJ6XOVE",
			HTMLURL: "https://twitchoncall.pagerduty.com/users/PJ6XOVE",
		},
		Name:     "Timur Kalandarov",
		Email:    "tikaland@justin.tv",
		Timezone: "America/New_York",
		Role:     "user",
		Teams: []pagerduty.Team{
			{
				APIObject: pagerduty.APIObject{
					ID:      "P83EOFI",
					Type:    "team_reference",
					Summary: "GSOC",
				},
			},
		},
	}
	return user
}

func getFakedSchedule() pagerduty.Schedule {
	schedule := pagerduty.Schedule{
		APIObject: pagerduty.APIObject{
			ID: "PUHMCXV",
		},
		Name: "",
	}
	return schedule
}

func getFakedEscalationPolicy() pagerduty.EscalationPolicy {
	policy := pagerduty.EscalationPolicy{}
	return policy
}

func getFakedOverridesList() pagerduty.ListOverridesResponse {
	overrides := pagerduty.ListOverridesResponse{
		Overrides: []pagerduty.Override{
			{
				ID:    "Q3WU06FHSCYOHG",
				Start: "2022-08-31T14:00:00-06:00",
				End:   "2022-09-01T00:00:00-06:00",
				User:  getFakedUser().APIObject,
			},
			{
				ID:    "Q1NF06I8X9HJAK",
				Start: "2022-09-01T14:00:00-06:00",
				End:   "2022-09-02T00:00:00-06:00",
				User:  getFakedUser().APIObject,
			},
			{
				ID:    "Q3BMMAACX0LQDM",
				Start: "2022-09-03T14:00:00-06:00",
				End:   "2022-09-04T00:00:00-06:00",
				User:  getFakedUser().APIObject,
			},
		},
	}
	return overrides
}

func getFakedTags() pagerduty.ListTagResponse {
	tags := pagerduty.ListTagResponse{
		Tags: []*pagerduty.Tag{
			{
				APIObject: pagerduty.APIObject{ID: "PRXFVK3", Summary: "ETS"},
				Label:     "ETS",
			},
			{
				APIObject: pagerduty.APIObject{ID: "P74RRGF", Summary: "GSOC"},
				Label:     "GSOC",
			},
			{
				APIObject: pagerduty.APIObject{ID: "PMRPRRZ", Summary: "VIDOPS"},
				Label:     "VIDOPS",
			},
			{
				APIObject: pagerduty.APIObject{ID: "P3V3R4S", Summary: "VOR"},
				Label:     "VOR",
			},
		},
	}
	return tags
}

func getFakedIncidents() []pagerduty.Incident {
	incidents := []pagerduty.Incident{
		{
			IncidentNumber: 2024046,
			Title:          "Twilight Automation Test Failure",
			Description:    "Twilight Automation Test Failure",
			CreatedAt:      "2022-09-06T03:00:15Z",
			Status:         "resolved",
			Urgency:        "low",
			APIObject:      pagerduty.APIObject{ID: "Q3XZW6AK6GZ3TZ"},
		},
		{
			IncidentNumber: 2024049,
			Title:          "Input Errors > 1000 over 1M",
			Description:    "cr01.sin04 tengige0/0/0/2/0 COLO:tm:cr01.bkk01:te0/0/0/1/0:US021-161:10G:::",
			CreatedAt:      "2022-09-06T03:04:51Z",
			Status:         "resolved",
			Urgency:        "low",
			APIObject:      pagerduty.APIObject{ID: "Q0FDXK79NN127I"},
		},
		{
			IncidentNumber: 2024046,
			Title:          "BatchGetCheckoutPrice primary p99 Latency > 0.8s",
			Description:    "BatchGetCheckoutPrice primary p99 Latency > 0.8s",
			CreatedAt:      "2022-09-06T03:18:04Z",
			Status:         "resolved",
			Urgency:        "low",
			APIObject:      pagerduty.APIObject{ID: "Q0HV1FERSUO36A"},
		},
	}

	return incidents
}

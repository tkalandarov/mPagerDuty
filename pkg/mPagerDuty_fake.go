package mPagerDuty

import (
	"fmt"
	"strings"
	"time"

	"github.com/PagerDuty/go-pagerduty"
)

type FakePDClient struct {
}

const (
	existingUserID       = "PJ6XOVE"
	existingUserName     = "Timur Kalandarov"
	existingScheduleID   = "PUHMCXV"
	existingScheduleName = "Caleb Young TESTING"
)

// GetOnCalls returns a pre-defined ListOnCallsResponse object defined in getFakedOnCalls()
// unless the array passed to the function contains no elements
func (fakeClient *FakePDClient) GetOnCallsByScheduleIDs(scheduleIDs []string) ([]pagerduty.OnCall, error) {
	if scheduleIDs == nil {
		return nil, fmt.Errorf("array of scheduleIDs must be defined")
	}

	if len(scheduleIDs) < 1 {
		return nil, fmt.Errorf("could not get on-calls because scheduleIDs was empty")
	}

	for i := 0; i < len(scheduleIDs); i++ {
		if len(strings.TrimSpace(scheduleIDs[i])) == 0 {
			return nil, fmt.Errorf("could not get on-calls because the passed array had empty string(s)")
		}
	}

	fakeResponse := getFakedOnCalls().OnCalls
	return fakeResponse, nil
}

// GetOnCalls always returns a pre-defined ListOnCallsResponse object defined in getFakedOnCalls()
func (fakeClient *FakePDClient) GetOnCallsWithOptions(options *pagerduty.ListOnCallOptions) (*pagerduty.ListOnCallsResponse, error) {
	fakeResponse := getFakedOnCalls()
	return fakeResponse, nil
}

// GetScheduleIDbyName returns a pre-defined same value if the passed argument is defined
// Otherwise, the function returns an error
func (fakeClient *FakePDClient) GetScheduleIDbyName(name string) (string, string, error) {
	if strings.TrimSpace(name) == "" {
		return "", "", fmt.Errorf("passed parameter 'name' must be specified")
	}
	return existingScheduleID, "America/Denver", nil
}

// GetUserIDbyName returns a pre-defined same value if the passed argument is defined
// Otherwise, the function returns an error
func (fakeClient *FakePDClient) GetUserIDbyName(name string) (string, error) {
	if strings.TrimSpace(name) == "" {
		return "", fmt.Errorf("passed parameter 'name' must be specified")
	}
	return existingUserID, nil
}

// GetUserByID returns a pre-defined pagerduty.User object if the passed 'id' parameter is defined
// Otherwise, the function returns an error
func (fakeClient *FakePDClient) GetUserByID(id string, options pagerduty.GetUserOptions) (*pagerduty.User, error) {
	if strings.TrimSpace(id) == "" {
		return nil, fmt.Errorf("passed parameter 'id' must be specified")
	}

	user := getFakedUser()
	return &user, nil
}

// ListAllUsers returns an array of three same pre-defined pagerduty.User objects
func (fakeClient *FakePDClient) ListAllUsers(options pagerduty.ListUsersOptions) ([]pagerduty.User, error) {
	users := []pagerduty.User{getFakedUser(), getFakedUser(), getFakedUser()}

	return users, nil
}

// GetOverrides returns a pagerduty.ListOverridesResponse containing an array
// of three different pre-defined pagerduty.Override objects, if all passed parameters are defined
// Otherwise, the function returns an error
func (fakeClient *FakePDClient) GetOverrides(scheduleID string, since string, until string, includeOverflow bool) (*pagerduty.ListOverridesResponse, error) {
	if strings.TrimSpace(scheduleID) == "" {
		return nil, fmt.Errorf("passed parameter 'scheduleID' must be specified")
	}

	if strings.TrimSpace(since) == "" || strings.TrimSpace(until) == "" {
		return nil, fmt.Errorf("passed parameters 'since' and 'until' dates must be specified")
	}

	overrides := getFakedOverridesList()
	return &overrides, nil
}

// CreateOverride will return nil and a faked *pagerduty.Override object
// if all passed parameters are defined
// Otherwise, the function returns an error
func (fakeClient *FakePDClient) CreateOverride(scheduleID string, userID string, start string, end string) (*pagerduty.Override, error) {
	if strings.TrimSpace(scheduleID) == "" {
		return nil, fmt.Errorf("passed parameter 'scheduleID' must be specified")
	}
	if strings.TrimSpace(userID) == "" {
		return nil, fmt.Errorf("passed parameter 'userID' must be specified")
	}
	if strings.TrimSpace(start) == "" {
		return nil, fmt.Errorf("passed parameter 'start' must be specified")
	}
	if strings.TrimSpace(end) == "" {
		return nil, fmt.Errorf("passed parameter 'end' must be specified")
	}
	return &pagerduty.Override{
		ID:    "TEST",
		Start: start,
		End:   end,
		User:  pagerduty.APIObject{ID: userID},
	}, nil
}

// RemoveOverride will return nil if all passed parameters are defined
// Otherwise, the function returns an error
func (fakeClient *FakePDClient) RemoveOverride(scheduleID string, overrideID string) error {
	if strings.TrimSpace(scheduleID) == "" {
		return fmt.Errorf("passed parameter 'scheduleID' must be specified")
	}
	if strings.TrimSpace(overrideID) == "" {
		return fmt.Errorf("passed parameter 'overrideID' must be specified")
	}
	return nil
}

// GetIndicentsByTag returns faked pre-defined array of pagerduty.Incident objects
func (fakeClient *FakePDClient) GetIndicentsByTag(tagName string, timeRange time.Duration) ([]pagerduty.Incident, error) {
	if strings.TrimSpace(tagName) == "" {
		return nil, fmt.Errorf("passed parameter 'tagName' must be specified")
	}

	incidentList := getFakedIncidents()
	return incidentList, nil
}

// CreateIncident will return a faked *pagerduty.Incident object
// if 'title' and 'serviceID' are defined.
// Otherwise, the function returns an error
func (fakeClient *FakePDClient) CreateIncident(title, serviceID, urgency, details, escalationPolicyID string) (*pagerduty.Incident, error) {
	// Only title and serviceID are required fields as per API reference
	if strings.TrimSpace(title) == "" {
		return nil, fmt.Errorf("passed parameter 'title' must be specified")
	}
	if strings.TrimSpace(serviceID) == "" {
		return nil, fmt.Errorf("passed parameter 'serviceID' must be specified")
	}
	return &pagerduty.Incident{
		Title:            title,
		Service:          pagerduty.APIObject{ID: serviceID},
		Urgency:          urgency,
		EscalationPolicy: pagerduty.APIObject{ID: escalationPolicyID},
	}, nil
}

// SearchIncidents will imitate searching and finding pagerduty.Incident objects
// by concatenating the value of serviceQuery with each element's Incident.Service.Summary property, if the parameter is defined.
// Otherwise, if serviceQuery is not defined, the function will return an empty array
func (fakeClient *FakePDClient) SearchIncidents(serviceQuery string, timeRange time.Duration) ([]pagerduty.Incident, error) {
	if serviceQuery == "" {
		return []pagerduty.Incident{}, nil
	}

	incidents := getFakedIncidents()
	for _, incident := range incidents {
		incident.Service.Summary += serviceQuery
	}
	return incidents, nil
}

// SearchIncidentLogs will return an error if the passed parameters are not specified and
// will return a faked pre-defined summary if the given logType is found in pre-defined list.
// Otherwise, the function will return nil
func (fakeClient *FakePDClient) SearchIncidentLogs(incidentID string, logType string) (*string, error) {
	if strings.TrimSpace(incidentID) == "" {
		return nil, fmt.Errorf("passed parameter 'incidentID' must be specified")
	}
	if strings.TrimSpace(logType) == "" {
		return nil, fmt.Errorf("passed parameter 'logType' must be specified")
	}

	if logType == "resolve_log_entry" {
		summary := "Resolved through the API."
		return &summary, nil
	}
	if logType == "assign_log_entry" {
		summary := "Assigned to Timur Kalandarov."
		return &summary, nil
	}
	if logType == "trigger_log_entry" {
		summary := "Triggered through the API."
		return &summary, nil
	}
	return nil, nil
}

// GetIndicentsByEscalationPolicy returns faked pre-defined array of pagerduty.Incident objects
func (fakeClient *FakePDClient) GetIndicentsByEscalationPolicy(escalationPolicyID string, timeRange time.Duration) ([]pagerduty.Incident, error) {
	if strings.TrimSpace(escalationPolicyID) == "" {
		return nil, fmt.Errorf("passed parameter 'escalationPolicyID' must be specified")
	}

	incidentList := getFakedIncidents()

	return incidentList, nil
}

// GetOnCallUsersIDs will return an array of N same pre-defied values
// where N is length of the array of names passed as a parameter
//
// If the parameter is nil or contains no elements, the function will return an error
func (fakeClient *FakePDClient) GetUsersIDsByNames(names []string) ([]string, error) {
	if names == nil {
		return nil, fmt.Errorf("array of names must be defined")
	}
	if len(names) < 1 {
		return nil, fmt.Errorf("array of names cannot be empty")
	}
	for i := 0; i < len(names); i++ {
		if len(strings.TrimSpace(names[i])) == 0 {
			return nil, fmt.Errorf("could not get user names because the passed array had empty string(s)")
		}
	}

	var resp []string
	for i := 0; i < len(names); i++ {
		resp = append(resp, existingUserID)
	}

	return resp, nil
}

// UpdateEscalationPolicy will return a faked *pagerduty.EscalationPolicy
// if 'id', 'userID', and 'escalation' parameters passed are defined.
// Otherwise, the function returns an error
func (fakeClient *FakePDClient) UpdateEscalationPolicy(id, userID, serviceID, teamID string, escalation []pagerduty.APIObject) (*pagerduty.EscalationPolicy, error) {
	// only id, userID, and escalation are required fields as per API reference
	if strings.TrimSpace(id) == "" {
		return nil, fmt.Errorf("passed parameter 'id' must be specified")
	}
	if strings.TrimSpace(userID) == "" {
		return nil, fmt.Errorf("passed parameter 'userID' must be specified")
	}
	if escalation == nil {
		return nil, fmt.Errorf("passed parameter 'escalation' cannot be nil")
	}

	return &pagerduty.EscalationPolicy{
		APIObject: pagerduty.APIObject{ID: id},
		Services:  []pagerduty.APIObject{{ID: serviceID}},
		Teams:     []pagerduty.APIReference{{ID: teamID}},
	}, nil
}

// GetEscalationPoliciesByTag returns a *pagerduty.ListEPResponse object
// containing an array of pre-defined escalation policies
func (fakeClient *FakePDClient) GetEscalationPoliciesByTag(tagID string) (*pagerduty.ListEPResponse, error) {
	if strings.TrimSpace(tagID) == "" {
		return nil, fmt.Errorf("passed parameter 'tagID' must be specified")
	}

	response := pagerduty.ListEPResponse{
		EscalationPolicies: []*pagerduty.APIObject{
			{
				ID:      "PP2PMMD",
				Summary: "Ad Server Escalation",
			},
			{
				ID:      "PKR3E6F",
				Summary: "App code validation",
			},
			{
				ID:      "PGPQHZF",
				Summary: "Browser Grid",
			},
		},
	}

	return &response, nil
}

// ListAllTags returns a *pagerduty.ListTagResponse object containing an array of pre-defined pagerduty.Tag objects
func (fakeClient *FakePDClient) ListAllTags(options pagerduty.ListTagOptions) ([]*pagerduty.Tag, error) {
	response := getFakedTags()
	return response.Tags, nil
}

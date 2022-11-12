package mPagerDuty

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/PagerDuty/go-pagerduty"
)

// MPagerDuty is a mockable interface that is used to deliver a Client object to dependent services
// and test without actually sending API requests to the Jira instance
type IMPagerDuty interface {
	// On-Calls
	GetOnCallsByScheduleIDs(scheduleIDs []string) ([]pagerduty.OnCall, error)
	GetOnCallsWithOptions(*pagerduty.ListOnCallOptions) (*pagerduty.ListOnCallsResponse, error)

	// Users
	ListAllUsers(options pagerduty.ListUsersOptions) ([]pagerduty.User, error)
	GetUserByID(id string, options pagerduty.GetUserOptions) (*pagerduty.User, error)
	GetUserIDbyName(name string) (string, error)
	GetUsersIDsByNames(names []string) ([]string, error)

	// Schedules
	GetScheduleIDbyName(name string) (string, string, error)

	// Overrides
	GetOverrides(scheduleID, since, until string, includeOverflow bool) (*pagerduty.ListOverridesResponse, error)
	CreateOverride(scheduleID string, userID string, start string, end string) (*pagerduty.Override, error)
	RemoveOverride(scheduleID string, overrideID string) error

	// Tags
	ListAllTags(options pagerduty.ListTagOptions) ([]*pagerduty.Tag, error)

	// Incidents
	GetIndicentsByEscalationPolicy(escalationPolicyID string, timeRange time.Duration) ([]pagerduty.Incident, error)
	GetIndicentsByTag(tagName string, timeRange time.Duration) ([]pagerduty.Incident, error)
	CreateIncident(title, serviceID, urgency, details, escalationPolicyID string) (*pagerduty.Incident, error)
	SearchIncidents(serviceQuery string, timeRange time.Duration) ([]pagerduty.Incident, error)
	SearchIncidentLogs(pdIncidentId string, logType string) (*string, error)

	// Escalation Policies
	GetEscalationPoliciesByTag(tagID string) (*pagerduty.ListEPResponse, error)
	UpdateEscalationPolicy(id, userID, serviceID, teamID string, escalation []pagerduty.APIObject) (*pagerduty.EscalationPolicy, error)
}

// GetMPagerDutyClient creates a usable actual go-pagerduty client or a usable faked client
// that returns dummy data when environment variables RUNNING_IN_JENKINS or LOCAL_DEV_TESTING
// are set to true
func GetMPagerDutyClient(authtoken string) (IMPagerDuty, error) {
	pdClient, err := newMPagerDutyClient(authtoken)
	if err != nil {
		return nil, err
	}
	return pdClient, nil
}

// GetOnCallsByScheduleIDs returns the list of all on-calls for the specified schedule IDs
// API documentation: https://developer.pagerduty.com/api-reference/3a6b910f11050-list-all-of-the-on-calls
func (c *client) GetOnCallsByScheduleIDs(scheduleIDs []string) ([]pagerduty.OnCall, error) {
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

	var onCalls []pagerduty.OnCall
	var offset uint = 0
	for {
		response, err := c.pdClient.ListOnCallsWithContext(context.Background(), pagerduty.ListOnCallOptions{
			TimeZone:    "UTC",
			ScheduleIDs: scheduleIDs,
			Limit:       limit,
			Total:       true,
			Offset:      offset,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to list on calls: %s", err)
		}

		onCalls = append(onCalls, response.OnCalls...)
		if len(response.OnCalls) < 1 {
			break
		}
		offset = offset + increase
	}

	return onCalls, nil
}

// GetOnCallsWithOptions returns the list of on-calls that satisfy the specified options query
func (c *client) GetOnCallsWithOptions(options *pagerduty.ListOnCallOptions) (*pagerduty.ListOnCallsResponse, error) {
	onCalls, error := c.pdClient.ListOnCallsWithContext(context.Background(), *options)
	return onCalls, error
}

// GetScheduleIDbyName returns ID and timezone for the schedule specified by name
// API documentation: https://developer.pagerduty.com/api-reference/3f03afb2c84a4-get-a-schedule
func (c *client) GetScheduleIDbyName(name string) (string, string, error) {
	if strings.TrimSpace(name) == "" {
		return "", "", fmt.Errorf("passed parameter 'name' must be specified")
	}

	var resp string = ""
	var tz string = ""
	var offset uint = 0

	for {
		schedules, err := c.pdClient.ListSchedulesWithContext(
			context.Background(),
			pagerduty.ListSchedulesOptions{
				Limit:  limit,
				Offset: offset,
				Query:  ""})
		if err != nil {
			return resp, tz, fmt.Errorf("error in getting schedules from PagerDuty: %s", err)
		}
		if len(schedules.Schedules) < 1 {
			break
		} else {
			for i := range schedules.Schedules {
				if strings.EqualFold(schedules.Schedules[i].Name, pdScheduleTitlePrefix+" "+name) {
					resp = schedules.Schedules[i].ID
					tz = schedules.Schedules[i].TimeZone
					break
				}
			}
			if resp != "" {
				break
			}
		}
		offset = offset + increase
	}
	return resp, tz, nil
}

// GetUserIDbyName returns ID for the user specified by name
func (c *client) GetUserIDbyName(name string) (string, error) {
	if strings.TrimSpace(name) == "" {
		return "", fmt.Errorf("passed parameter 'name' must be specified")
	}

	var id string = ""
	var offset uint = 0

	for {
		users, err := c.pdClient.ListUsersWithContext(
			context.Background(),
			pagerduty.ListUsersOptions{
				Limit:   limit,
				Offset:  offset,
				TeamIDs: pgGSOCTeamID})
		if err != nil {
			return "", fmt.Errorf("error in getting userID from PagerDuty: %s", err)
		}

		if len(users.Users) < 1 {
			break
		} else {
			for i := range users.Users {
				normalizedName, err := normalizeString(users.Users[i].Name)
				if err != nil {
					return "", fmt.Errorf("error while normalizing user name '%s': %s", users.Users[i].Name, err)
				}

				if strings.EqualFold(normalizedName, name) {
					id = users.Users[i].ID
					break
				}
			}
			if id != "" {
				break
			}
		}
		offset = offset + increase
	}
	return id, nil
}

// GetUserByID returns *pagerduty.User object for the user specified by ID
// API documentation: https://developer.pagerduty.com/api-reference/2395ca1feb25e-get-a-user
func (c *client) GetUserByID(id string, options pagerduty.GetUserOptions) (*pagerduty.User, error) {
	if strings.TrimSpace(id) == "" {
		return nil, fmt.Errorf("passed parameter 'id' must be specified")
	}

	user, err := c.pdClient.GetUserWithContext(
		context.Background(),
		id, options)

	if err != nil {
		return nil, fmt.Errorf("failed to get pagerduty user: %s", err)
	}

	return user, nil
}

// ListAllUsers returns *pagerduty.ListUsersResponse object that contains infotmation about all PagerDuty users
// API documentation: https://developer.pagerduty.com/api-reference/c96e889522dd6-list-users
func (c *client) ListAllUsers(options pagerduty.ListUsersOptions) ([]pagerduty.User, error) {
	options.Limit = limit
	options.Offset = 0

	var users []pagerduty.User
	var offset uint = 0
	for {
		response, err := c.pdClient.ListUsersWithContext(context.Background(), options)
		if err != nil {
			return nil, fmt.Errorf("failed to list users: %s", err)
		}

		users = append(users, response.Users...)
		if len(response.Users) < 1 {
			break
		}
		offset = offset + increase
	}

	return users, nil
}

// GetOverrides returns *pagerduty.ListOverridesResponse object for the specified name, schedule ID, start and end dates
// API documentation: https://developer.pagerduty.com/api-reference/cb747199f63a9-list-overrides
func (c *client) GetOverrides(scheduleID, since, until string, includeOverflow bool) (*pagerduty.ListOverridesResponse, error) {
	if strings.TrimSpace(scheduleID) == "" {
		return nil, fmt.Errorf("passed parameter 'scheduleID' must be specified")
	}

	if strings.TrimSpace(since) == "" || strings.TrimSpace(until) == "" {
		return nil, fmt.Errorf("passed parameters 'since' and 'until' dates must be specified")
	}

	overrides, err := c.pdClient.ListOverridesWithContext(
		context.Background(),
		scheduleID, pagerduty.ListOverridesOptions{
			Since:    since,
			Until:    until,
			Overflow: includeOverflow, // unless this parameter is set to true, any entry that passes the date range bounds will be truncated at the bounds
		})
	if err != nil {
		return nil, err
	}
	return overrides, nil
}

// CreateOverride creates override for the specified by ID schedule and user, as well as start and end dates,
// and returns a newly created *pagerduty.Override object in case of success, and error, in case of failure
func (c *client) CreateOverride(scheduleID string, userID string, start string, end string) (*pagerduty.Override, error) {
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

	newOverride, err := c.pdClient.CreateOverrideWithContext(
		context.Background(),
		scheduleID, pagerduty.Override{Start: start, End: end, User: pagerduty.APIObject{ID: userID, Type: "user"}})
	if err != nil {
		return nil, fmt.Errorf("error while creating override on PagerDuty: %s", err)
	}
	return newOverride, nil
}

// RemoveOverride deletes override specified by its ID and schedule ID
func (c *client) RemoveOverride(scheduleID string, overrideID string) error {
	if strings.TrimSpace(scheduleID) == "" {
		return fmt.Errorf("passed parameter 'scheduleID' must be specified")
	}
	if strings.TrimSpace(overrideID) == "" {
		return fmt.Errorf("passed parameter 'overrideID' must be specified")
	}

	err := c.pdClient.DeleteOverrideWithContext(context.Background(), scheduleID, overrideID)
	if err != nil {
		return fmt.Errorf("error while removing override on PagerDuty: %s", err)
	}
	return nil
}

// GetIndicentsByEscalationPolicy returns an array of pagerduty.Incident objects
// assosiated with the specified escalation policy within the given time range
//
// API documentation: https://developer.pagerduty.com/api-reference/5a579467410f7-get-connected-entities
func (c *client) GetIndicentsByEscalationPolicy(escalationPolicyID string, timeRange time.Duration) ([]pagerduty.Incident, error) {
	if strings.TrimSpace(escalationPolicyID) == "" {
		return nil, fmt.Errorf("passed parameter 'escalationPolicyID' must be specified")
	}

	UTC, err := time.LoadLocation("UTC")
	if err != nil {
		return nil, err
	}
	var since = ""
	querySinceTime := time.Now().Add(timeRange * time.Minute).In(UTC).Format("2006-01-02T15:04:05")
	since = querySinceTime

	var incidentList []pagerduty.Incident
	var offset uint = 0
	for {
		var options = pagerduty.ListIncidentsOptions{
			Limit:    limit,
			Offset:   offset,
			Since:    since,
			TimeZone: "UTC",
		}
		response, err := c.pdClient.ListIncidentsWithContext(context.Background(), options)
		if err != nil {
			return nil, fmt.Errorf("failed to list incidents: %s", err)
		}

		if len(response.Incidents) < 1 {
			break
		}
		for _, incident := range response.Incidents {
			if incident.EscalationPolicy.ID == escalationPolicyID {
				incidentList = append(incidentList, incident)
			}
			since = incident.CreatedAt
		}

		offset = offset + increase
	}
	return incidentList, nil
}

// GetIndicentsByTag returns an array of pagerduty.Incident objects
// assosiated with the specified tag within the given time range
func (c *client) GetIndicentsByTag(tagName string, timeRange time.Duration) ([]pagerduty.Incident, error) {
	if strings.TrimSpace(tagName) == "" {
		return nil, fmt.Errorf("passed parameter 'tagName' must be specified")
	}

	var incidentList []pagerduty.Incident
	tags, err := c.ListAllTags(pagerduty.ListTagOptions{Query: tagName})
	if err != nil {
		return nil, err
	}
	tagIdList := tags

	epResponse, err := c.GetEscalationPoliciesByTag(tagIdList[0].ID)
	if err != nil {
		return nil, err
	}
	epList := epResponse.EscalationPolicies

	for _, ep := range epList {
		incidents, err := c.GetIndicentsByEscalationPolicy(ep.ID, timeRange)
		if err != nil {
			return nil, err
		}
		incidentList = append(incidentList, incidents...)
	}
	return incidentList, nil
}

// CreateIncident creates PagerDuty incident configured with input parameters
func (c *client) CreateIncident(title, serviceID, urgency, details, escalationPolicyID string) (*pagerduty.Incident, error) {
	// Only title and serviceID are required fields as per API reference
	if strings.TrimSpace(title) == "" {
		return nil, fmt.Errorf("passed parameter 'title' must be specified")
	}
	if strings.TrimSpace(serviceID) == "" {
		return nil, fmt.Errorf("passed parameter 'serviceID' must be specified")
	}

	options := &pagerduty.CreateIncidentOptions{
		Type:  "incident",
		Title: title,
		Service: &pagerduty.APIReference{
			ID:   serviceID,
			Type: "service_reference",
		},
	}

	if len(urgency) > 0 {
		options.Urgency = urgency
	}
	if len(details) > 0 {
		options.Body = &pagerduty.APIDetails{
			Type:    "incident_body",
			Details: details,
		}
	}
	if len(escalationPolicyID) > 0 {
		options.EscalationPolicy = &pagerduty.APIReference{
			ID:   escalationPolicyID,
			Type: "escalation_policy_reference",
		}
	}

	incident, err := c.pdClient.CreateIncidentWithContext(context.Background(), "nobody@justin.tv", options)
	if err != nil {
		return nil, fmt.Errorf("API request to create an incident failed: %s", err)
	}

	return incident, nil
}

// SearchIncidents returns an array of pagerduty.Incident objects,
// whose Incident.Service.Summary property contains a specified serviceQuery
func (c *client) SearchIncidents(serviceQuery string, timeRange time.Duration) ([]pagerduty.Incident, error) {
	UTC, err := time.LoadLocation("UTC")
	if err != nil {
		return nil, fmt.Errorf("failed to load time location: %s", err)
	}

	queryTimeRange := time.Now().Add(timeRange * time.Minute).In(UTC).Format("2006-01-02T15:04:05")
	var since = ""
	since = queryTimeRange

	var searchResults []pagerduty.Incident
	var offset uint = 0
	for {
		var options = pagerduty.ListIncidentsOptions{
			Limit:    limit,
			Offset:   offset,
			Since:    since,
			TimeZone: "UTC",
		}
		response, err := c.pdClient.ListIncidentsWithContext(context.Background(), options)
		if err != nil {
			return nil, fmt.Errorf("failed to list incidents: %s", err)
		}

		if len(response.Incidents) < 1 {
			break
		}
		for _, incident := range response.Incidents {
			if strings.Contains(incident.Service.Summary, serviceQuery) {
				searchResults = append(searchResults, incident)
			}
			since = incident.CreatedAt
		}

		offset = offset + increase
	}

	return searchResults, nil
}

// SearchIncidentLogs will return a nilable string that contains entries of
// the given type that are asossiated with the specified incidentID,
// or nil otherwise.
func (c *client) SearchIncidentLogs(incidentID string, logType string) (*string, error) {
	if strings.TrimSpace(incidentID) == "" {
		return nil, fmt.Errorf("passed parameter 'incidentID' must be specified")
	}
	if strings.TrimSpace(logType) == "" {
		return nil, fmt.Errorf("passed parameter 'logType' must be specified")
	}

	var offset uint = 0
	for {
		var options = pagerduty.ListIncidentLogEntriesOptions{
			Limit:      limit,
			Offset:     offset,
			IsOverview: false,
		}
		// API Documentation: https://developer.pagerduty.com/api-reference/367602cbc1c28-list-log-entries-for-an-incident
		incidentLog, err := c.pdClient.ListIncidentLogEntriesWithContext(
			context.Background(), incidentID, options)
		if err != nil {
			return nil, fmt.Errorf("error getting incident log: %s", err)
		}

		if len(incidentLog.LogEntries) < 1 {
			break
		}

		for _, logEntry := range incidentLog.LogEntries {
			if logEntry.Type == logType {
				return &logEntry.Summary, nil
			}
		}
	}

	return nil, nil
}

// GetUsersIDsByNames returns an array of user IDs that match specified names
func (c *client) GetUsersIDsByNames(names []string) ([]string, error) {
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
	for i := range names {
		users, err := c.ListAllUsers(
			pagerduty.ListUsersOptions{
				Query: names[i],
			})
		if err != nil {
			return []string{}, fmt.Errorf("error getting users from PagerDuty: %s", err)
		}

		for j := range users {
			if strings.EqualFold(users[j].Name, names[i]) {
				resp = append(resp, users[j].ID)
				break
			}
		}
		if len(resp) > 0 {
			break
		}
	}
	return resp, nil
}

// UpdateEscalationPolicy updates escalation policy using the input parameters
func (c *client) UpdateEscalationPolicy(id, userID, serviceID, teamID string, escalation []pagerduty.APIObject) (*pagerduty.EscalationPolicy, error) {
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

	escalationPolicy := pagerduty.EscalationPolicy{
		EscalationRules: []pagerduty.EscalationRule{
			{
				Delay: 5,
				Targets: []pagerduty.APIObject{
					{
						ID:   userID,
						Type: "user_reference",
					},
				},
			},
			{
				Delay:   15,
				Targets: escalation,
			},
		},
	}

	if len(serviceID) > 0 {
		escalationPolicy.Services = []pagerduty.APIObject{
			{
				ID:   serviceID,
				Type: "service_reference",
			},
		}
	}

	if len(teamID) > 0 {
		escalationPolicy.Teams = []pagerduty.APIReference{
			{
				ID:   teamID,
				Type: "team_reference",
			},
		}
	}

	newEscalationPolicy, err := c.pdClient.UpdateEscalationPolicyWithContext(context.Background(), id, escalationPolicy)
	return newEscalationPolicy, err
}

// GetEscalationPoliciesByTag returns a *pagerduty.ListEPResponse object containing an array of all existing escalation policies
func (c *client) GetEscalationPoliciesByTag(tagID string) (*pagerduty.ListEPResponse, error) {
	if strings.TrimSpace(tagID) == "" {
		return nil, fmt.Errorf("passed parameter 'tagID' must be specified")
	}

	response, err := c.pdClient.GetEscalationPoliciesByTag(tagID)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// ListAllTags returns a *pagerduty.ListTagResponse containing all existing tags
// API documentation: https://developer.pagerduty.com/api-reference/e44b160c69bf3-list-tags
func (c *client) ListAllTags(options pagerduty.ListTagOptions) ([]*pagerduty.Tag, error) {
	options.Limit = limit
	options.Offset = 0

	var tags []*pagerduty.Tag
	var offset uint = 0
	for {
		response, err := c.pdClient.ListTags(options)
		if err != nil {
			return nil, fmt.Errorf("failed to list tags: %s", err)
		}

		tags = append(tags, response.Tags...)
		if len(response.Tags) < 1 {
			break
		}
		offset = offset + increase
	}

	return tags, nil
}

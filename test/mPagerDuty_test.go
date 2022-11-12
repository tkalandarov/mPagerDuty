package mPagerDuty_test

import (
	"os"
	"testing"
	"time"

	mPagerDuty "mpagerduty/pkg"

	"github.com/PagerDuty/go-pagerduty"
	"github.com/stretchr/testify/assert"
)

const authtoken string = "AUTHTOKEN"

/* THE TESTS BELOW THIS COMMENT ARE UNIT TESTS THAT WILL RUN WITHOUT SENDING
   API REQUEST OUT ACROSS THE NETWORK. THESE SHOULD BE QUICK AND WILL RUN
   DURING NORMAL JENKINS BUILDS */

func TestGetOnCallsByScheduleIDs(t *testing.T) {
	tests := []struct {
		scheduleIDs []string
		expectedErr bool
	}{
		{
			[]string{"PUY4P9O", "P10QVCS"},
			false,
		},
		{
			[]string{"PUY4P9O"},
			false,
		},
		{
			[]string{"P10QVCS", ""},
			true,
		},
		{
			nil,
			true,
		},
		{
			[]string{},
			true,
		},
	}

	mPD := mPagerDuty.FakePDClient{}
	for _, test := range tests {
		_, err := mPD.GetOnCallsByScheduleIDs(test.scheduleIDs)
		if test.expectedErr {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err)
		}
	}
}

func TestGetScheduleIDbyName(t *testing.T) {
	tests := []struct {
		name        string
		expectedErr bool
	}{
		{
			"Timur Kalandarov",
			false,
		},
		{
			"    ",
			true,
		},
	}

	mPD := mPagerDuty.FakePDClient{}
	for _, test := range tests {
		_, _, err := mPD.GetScheduleIDbyName(test.name)
		if test.expectedErr {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err)
		}
	}
}

func TestGetUserIDbyName(t *testing.T) {
	tests := []struct {
		name           string
		expectedResult string
		expectedErr    bool
	}{
		{
			"Timur Kalandarov",
			"PUHMCXV",
			false,
		},
		{
			"    ",
			"",
			true,
		},
	}

	mPD := mPagerDuty.FakePDClient{}
	for _, test := range tests {
		result, _, err := mPD.GetScheduleIDbyName(test.name)
		if test.expectedErr {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err)
			assert.Equal(t, result, test.expectedResult)
		}
	}
}

func TestGetUserByID(t *testing.T) {
	tests := []struct {
		userID           string
		expectedUserName string
		expectedErr      bool
	}{
		{
			"PUHMCXV",
			"Timur Kalandarov",
			false,
		},
		{
			"     ",
			"",
			true,
		},
		{
			"    TEST",
			"Timur Kalandarov",
			false,
		},
	}

	mPD := mPagerDuty.FakePDClient{}
	for _, test := range tests {
		user, err := mPD.GetUserByID(test.userID, pagerduty.GetUserOptions{})
		if test.expectedErr {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err)
			assert.Equal(t, user.Name, test.expectedUserName)
		}
	}
}

func TestGetOverrides(t *testing.T) {
	var minTime = time.Now().Format(time.RFC3339)
	var maxTime = time.Now().AddDate(0, 2, 0).Format(time.RFC3339)

	tests := []struct {
		scheduleID  string
		userName    string
		expectedErr bool
	}{
		{
			"P10QVCS",
			"Timur Kalandarov",
			false,
		},
		{
			"   P10QVCS",
			"Timur Kalandarov",
			false,
		},
		{
			"    ",
			"",
			true,
		},
	}

	mPD := mPagerDuty.FakePDClient{}
	for _, test := range tests {
		resp, err := mPD.GetOverrides(test.scheduleID, minTime, maxTime, false)
		if test.expectedErr {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err)
			assert.NotEmpty(t, resp.Overrides)
		}
	}
}

func TestCreateOverride(t *testing.T) {
	var minTime = time.Now().Format(time.RFC3339)
	var maxTime = time.Now().AddDate(0, 2, 0).Format(time.RFC3339)

	tests := []struct {
		scheduleID  string
		userName    string
		expectedErr bool
	}{
		{
			"P10QVCS",
			"Timur Kalandarov",
			false,
		},
		{
			"  ",
			"",
			true,
		},
		{
			"    TEST",
			"",
			true,
		},
		{
			"",
			"TETS  ",
			true,
		},
	}

	mPD := mPagerDuty.FakePDClient{}
	for _, test := range tests {
		override, err := mPD.CreateOverride(test.userName, test.scheduleID, minTime, maxTime)
		if test.expectedErr {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err)
			assert.NotNil(t, override)
		}
	}
}

func TestRemoveOverride(t *testing.T) {
	tests := []struct {
		scheduleID  string
		overrideID  string
		expectedErr bool
	}{
		{
			"P10QVCS",
			"Timur Kalandarov",
			false,
		},
		{
			" TEST",
			"",
			true,
		},
		{
			"  ",
			"  TEST",
			true,
		},
	}

	mPD := mPagerDuty.FakePDClient{}
	for _, test := range tests {
		err := mPD.RemoveOverride(test.scheduleID, test.overrideID)
		if test.expectedErr {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err)
		}
	}
}

func TestGetIndicentsByEP(t *testing.T) {
	tests := []struct {
		escalationPolicyID string
		expectedErr        bool
	}{
		{
			"PXYKJ4K",
			false,
		},
		{
			"",
			true,
		},
		{
			"   ",
			true,
		},
	}

	mPD := mPagerDuty.FakePDClient{}
	for _, test := range tests {
		incidents, err := mPD.GetIndicentsByEscalationPolicy(test.escalationPolicyID, -30)
		if test.expectedErr {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err)
			assert.NotEmpty(t, incidents)
		}
	}
}

func TestGetIndicentsByTag(t *testing.T) {
	tests := []struct {
		tagID       string
		expectedErr bool
	}{
		{
			"P74RRGF", // GSOC tag
			false,
		},
		{
			"",
			true,
		},
		{
			"   ",
			true,
		},
	}

	mPD := mPagerDuty.FakePDClient{}
	for _, test := range tests {
		incidents, err := mPD.GetIndicentsByTag(test.tagID, -30)
		if test.expectedErr {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err)
			assert.NotEmpty(t, incidents)
		}
	}
}

func TestCreateIncident(t *testing.T) {
	tests := []struct {
		title              string
		serviceID          string
		urgency            string
		details            string
		escalationPolicyID string
		expectedErr        bool
	}{
		{
			"The server is on fire",
			"P03NRF0",
			"high",
			"Someone call firefighers!",
			"PXYKJ4K",
			false,
		},
		{
			"",
			"P03NRF0",
			"high",
			"Title is a requried field. This test should throw an error.",
			"PXYKJ4K",
			true,
		},
		{
			"The server is on fire",
			"",
			"high",
			"ServiceID is also required. This one should fail too.",
			"PXYKJ4K",
			true,
		},
		{
			"The server is on fire",
			"P03NRF0",
			"",
			"",
			"",
			false,
		},
	}

	mPD := mPagerDuty.FakePDClient{}
	for _, test := range tests {
		incident, err := mPD.CreateIncident(test.title, test.serviceID, test.urgency, test.details, test.escalationPolicyID)
		if test.expectedErr {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err)
			assert.NotNil(t, incident)
			assert.Equal(t, incident.Title, test.title)
		}
	}
}

func TestGetUsersIDsByNames(t *testing.T) {
	tests := []struct {
		userNames   []string
		expectedErr bool
	}{
		{
			[]string{"Timur Kalandarov"}, // PWKNFGT
			false,
		},
		{
			[]string{"Timur Kalandarov", "   "},
			true,
		},
		{
			nil,
			true,
		},
	}

	mPD := mPagerDuty.FakePDClient{}
	for _, test := range tests {
		users, err := mPD.GetUsersIDsByNames(test.userNames)
		if test.expectedErr {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err)
			assert.Equal(t, len(users), len(test.userNames))
		}
	}
}

func TestGetEscalationPoliciesByTag(t *testing.T) {
	tests := []struct {
		tagID       string
		expectedErr bool
	}{
		{
			"P74RRGF", // GSOC tag
			false,
		},
		{
			"",
			true,
		},
		{
			"   ",
			true,
		},
	}

	mPD := mPagerDuty.FakePDClient{}
	for _, test := range tests {
		_, err := mPD.GetEscalationPoliciesByTag(test.tagID)
		if test.expectedErr {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err)
		}
	}
}

func TestUpdateEscalationPolicy(t *testing.T) {
	tests := []struct {
		id          string
		userID      string
		serviceID   string
		teamID      string
		escalation  []pagerduty.APIObject
		expectedErr bool
	}{
		{
			"P23N6LT", // INCY DEV Escalation Policy
			"PWKNFGT", // Timur Kalandarov
			"P39MEWZ", // GSOC General Service
			"P83EOFI", // GSOC Team
			[]pagerduty.APIObject{
				{
					ID:   "P273W1N", // Caleb Young
					Type: "user_reference",
				},
			},
			false,
		},
		{
			"   P23N6LT", // INCY DEV Escalation Policy
			"PWKNFGT",    // Timur Kalandarov
			"  ",
			"",
			[]pagerduty.APIObject{
				{
					ID:   "P273W1N", // Caleb Young
					Type: "user_reference",
				},
			},
			false,
		},
		{
			"",
			"TEST",
			"TEST",
			"TEST",
			[]pagerduty.APIObject{
				{
					ID:   "TEST",
					Type: "user_reference",
				},
			},
			true,
		},
		{
			"TEST",
			"",
			"TEST",
			"TEST",
			[]pagerduty.APIObject{
				{
					ID:   "TEST",
					Type: "user_reference",
				},
			},
			true,
		},
		{
			"TEST",
			"  TEST",
			"TEST",
			"TEST",
			nil,
			true,
		},
	}

	mPD := mPagerDuty.FakePDClient{}
	for _, test := range tests {
		escalationPolicy, err := mPD.UpdateEscalationPolicy(test.id, test.userID, test.serviceID, test.teamID, test.escalation)
		if test.expectedErr {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err)
			assert.Equal(t, escalationPolicy.ID, test.id)
		}
	}
}

/* THE TESTS BELOW THIS COMMENT ARE FULL INTEGRATION TESTS AND WILL USE THE
   ACTUAL REAL TWITCH'S PAGERDUTY INSTANCE. BECAUSE OF THIS, THESE TESTS
   SHOULD NOT BE RUN FREQUENTLY AND HAVE A LINE AT THE BEGINNING OF THEM
   SO THEY ARE SKIPPED DURING THE JENKINS BUILD PROCESS. HOWEVER, THESE TESTS
   SHOULD BE RUN EVERY SO OFTEN LOCALLY TO VERIFY THEY STILL WORK CORRECTLY */

func TestGetOnCallsByScheduleIDsIntegration(t *testing.T) {
	if os.Getenv("RUNNING_IN_JENKINS") == "true" {
		t.Skip("Skipping mPagerDuty TestGetOnCallsByScheduleIDsIntegration in CI")
	}

	tests := []struct {
		scheduleIDs []string
		expectedErr bool
	}{
		{
			[]string{"PUY4P9O"},
			false,
		},
		{
			[]string{"P10QVCS"},
			false,
		},
		{
			[]string{"DUMMYSTRING"},
			true,
		},
	}

	mPD, _ := mPagerDuty.GetMPagerDutyClient(authtoken)
	for _, test := range tests {
		_, err := mPD.GetOnCallsByScheduleIDs(test.scheduleIDs)
		if test.expectedErr {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err)
		}
	}
}

func TestGetScheduleIDbyNameIntegration(t *testing.T) {
	if os.Getenv("RUNNING_IN_JENKINS") == "true" {
		t.Skip("Skipping mPagerDuty TestGetScheduleIDbyNameIntegration in CI")
	}

	tests := []struct {
		name        string
		expectedErr bool
	}{
		{
			"Timur Kalandarov",
			false,
		},
		{
			"DUMMYSTRING",
			true,
		},
	}

	mPD, _ := mPagerDuty.GetMPagerDutyClient(authtoken)
	for _, test := range tests {
		_, _, err := mPD.GetScheduleIDbyName(test.name)
		if test.expectedErr {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err)
		}
	}
}

func TestGetUserIDbyNameIntegration(t *testing.T) {
	if os.Getenv("RUNNING_IN_JENKINS") == "true" {
		t.Skip("Skipping mPagerDuty TestGetUserIDbyNameIntegration in CI")
	}

	tests := []struct {
		name           string
		expectedResult string
		expectedErr    bool
	}{
		{
			"Timur Kalandarov",
			"PWKNFGT",
			false,
		},
		{
			"DUMMYSTRING",
			"",
			true,
		},
	}

	mPD, _ := mPagerDuty.GetMPagerDutyClient(authtoken)
	for _, test := range tests {
		result, _, err := mPD.GetScheduleIDbyName(test.name)
		if test.expectedErr {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err)
			assert.Equal(t, result, test.expectedResult)
		}
	}
}

func TestGetUserByIDIntegration(t *testing.T) {
	if os.Getenv("RUNNING_IN_JENKINS") == "true" {
		t.Skip("Skipping mPagerDuty TestGetUserByIDIntegration in CI")
	}

	tests := []struct {
		userID           string
		expectedUserName string
		expectedErr      bool
	}{
		{
			"PWKNFGT",
			"Timur Kalandarov",
			false,
		},
		{
			"DUMMYSTRING",
			"",
			true,
		},
		{
			"    TEST",
			"",
			true,
		},
	}

	mPD, _ := mPagerDuty.GetMPagerDutyClient(authtoken)
	for _, test := range tests {
		user, err := mPD.GetUserByID(test.userID, pagerduty.GetUserOptions{})
		if test.expectedErr {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err)
			assert.Equal(t, user.Name, test.expectedUserName)
		}
	}
}

func TestListAllUsersIntegration(t *testing.T) {
	if os.Getenv("RUNNING_IN_JENKINS") == "true" {
		t.Skip("Skipping mPagerDuty TestListAllUsersIntegration in CI")
	}

	mPD, _ := mPagerDuty.GetMPagerDutyClient(authtoken)

	response, err := mPD.ListAllUsers(pagerduty.ListUsersOptions{})

	assert.Nil(t, err)
	assert.NotNil(t, response)
}

func TestGetOverridesIntegration(t *testing.T) {
	if os.Getenv("RUNNING_IN_JENKINS") == "true" {
		t.Skip("Skipping mPagerDuty TestGetOverridesIntegration in CI")
	}

	var minTime = time.Now().Format(time.RFC3339)
	var maxTime = time.Now().AddDate(0, 2, 0).Format(time.RFC3339)

	tests := []struct {
		scheduleID  string
		userName    string
		expectedErr bool
	}{
		{
			"P10QVCS",
			"Timur Kalandarov",
			false,
		},
		{
			"DUMMYSTRING",
			"",
			true,
		},
		{
			"    TEST",
			"",
			true,
		},
	}

	mPD, _ := mPagerDuty.GetMPagerDutyClient(authtoken)
	for _, test := range tests {
		resp, err := mPD.GetOverrides(test.scheduleID, minTime, maxTime, false)
		if test.expectedErr {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err)
			assert.NotEmpty(t, resp.Overrides)
		}
	}
}

func TestCreateOverrideIntegration(t *testing.T) {
	if os.Getenv("RUNNING_IN_JENKINS") == "true" {
		t.Skip("Skipping mPagerDuty TestCreateOverrideIntegration in CI")
	}

	var minTime = time.Now().Format(time.RFC3339)
	var maxTime = time.Now().AddDate(0, 2, 0).Format(time.RFC3339)

	tests := []struct {
		scheduleID  string
		userName    string
		expectedErr bool
	}{
		{
			"P10QVCS",
			"Timur Kalandarov",
			false,
		},
		{
			"DUMMYSTRING",
			"",
			true,
		},
		{
			"    TEST",
			"",
			true,
		},
	}

	mPD, _ := mPagerDuty.GetMPagerDutyClient(authtoken)
	for _, test := range tests {
		override, err := mPD.CreateOverride(test.userName, test.scheduleID, minTime, maxTime)
		if test.expectedErr {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err)
			assert.NotNil(t, override)
		}
	}
}

func TestRemoveOverrideIntegration(t *testing.T) {
	if os.Getenv("RUNNING_IN_JENKINS") == "true" {
		t.Skip("Skipping mPagerDuty TestRemoveOverrideIntegration in CI")
	}

	tests := []struct {
		scheduleID  string
		userName    string
		expectedErr bool
	}{
		{
			"P10QVCS",
			"Timur Kalandarov",
			false,
		},
		{
			"DUMMYSTRING",
			"",
			true,
		},
		{
			"    TEST",
			"",
			true,
		},
	}

	mPD, _ := mPagerDuty.GetMPagerDutyClient(authtoken)
	for _, test := range tests {
		err := mPD.RemoveOverride(test.userName, test.scheduleID)
		if test.expectedErr {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err)
		}
	}
}

func TestGetIndicentsByEPIntegration(t *testing.T) {
	if os.Getenv("RUNNING_IN_JENKINS") == "true" {
		t.Skip("Skipping mPagerDuty TestGetIndicentsByEPIntegration in CI")
	}

	tests := []struct {
		escalationPolicyID string
		expectedErr        bool
	}{
		{
			"P23N6LT", // INCY DEV Escalation Policy
			false,
		},
		{
			"",
			true,
		},
		{
			"   ",
			true,
		},
		{
			"TEST",
			true,
		},
	}

	mPD, _ := mPagerDuty.GetMPagerDutyClient(authtoken)
	for _, test := range tests {
		_, err := mPD.GetIndicentsByEscalationPolicy(test.escalationPolicyID, -30)
		if test.expectedErr {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err)
		}
	}
}

func TestGetIndicentsByTagIntegration(t *testing.T) {
	if os.Getenv("RUNNING_IN_JENKINS") == "true" {
		t.Skip("Skipping mPagerDuty TestGetIndicentsByTagIntegration in CI")
	}

	tests := []struct {
		tagID       string
		expectedErr bool
	}{
		{
			"P74RRGF", // GSOC tag
			false,
		},
		{
			"",
			true,
		},
		{
			"   ",
			true,
		},
		{
			"TEST",
			true,
		},
	}

	mPD, _ := mPagerDuty.GetMPagerDutyClient(authtoken)
	for _, test := range tests {
		_, err := mPD.GetIndicentsByTag(test.tagID, -30)
		if test.expectedErr {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err)
		}
	}
}

func TestListAllTagsIntegration(t *testing.T) {
	if os.Getenv("RUNNING_IN_JENKINS") == "true" {
		t.Skip("Skipping mPagerDuty TestListAllTagsIntegration in CI")
	}

	mPD, _ := mPagerDuty.GetMPagerDutyClient(authtoken)

	tags, err := mPD.ListAllTags(pagerduty.ListTagOptions{})

	assert.Nil(t, err)
	assert.NotNil(t, tags)
	assert.Greater(t, len(tags), 0)
}

func TestCreateIncidentIntegration(t *testing.T) {
	if os.Getenv("RUNNING_IN_JENKINS") == "true" {
		t.Skip("Skipping mPagerDuty TestCreateIncidentIntegration in CI")
	}

	tests := []struct {
		title              string
		serviceID          string
		urgency            string
		details            string
		escalationPolicyID string
		expectedErr        bool
	}{
		{
			"The server is on fire",
			"P03NRF0",
			"high",
			"Someone call firefighers!",
			"PXYKJ4K",
			false,
		},
		{
			"",
			"P03NRF0",
			"high",
			"Title is a requried field. This test should throw an error.",
			"PXYKJ4K",
			true,
		},
		{
			"The server is on fire",
			"",
			"high",
			"ServiceID is also required. This one should fail too.",
			"PXYKJ4K",
			true,
		},
		{
			"The server is on fire",
			"P03NRF0",
			"",
			"",
			"",
			false,
		},
	}

	mPD, _ := mPagerDuty.GetMPagerDutyClient(authtoken)
	for _, test := range tests {
		incident, err := mPD.CreateIncident(test.title, test.serviceID, test.urgency, test.details, test.escalationPolicyID)
		if test.expectedErr {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err)
			assert.NotNil(t, incident)
			assert.Equal(t, incident.Title, test.title)
		}
	}
}

func TestGetUsersIDsByNamesIntegration(t *testing.T) {
	if os.Getenv("RUNNING_IN_JENKINS") == "true" {
		t.Skip("Skipping mPagerDuty TestGetOnCallUsersIDsIntegration in CI")
	}

	tests := []struct {
		userNames   []string
		expectedErr bool
	}{
		{
			[]string{"Timur Kalandarov"}, // PWKNFGT
			false,
		},
		{
			[]string{"Timur Kalandarov", "Caleb Young"}, // PWKNFGT, P273W1N
			false,
		},
		{
			[]string{"DUMMYSTRING"},
			true,
		},
	}

	mPD, _ := mPagerDuty.GetMPagerDutyClient(authtoken)
	for _, test := range tests {
		users, _ := mPD.GetUsersIDsByNames(test.userNames)
		if test.expectedErr {
			assert.NotEqual(t, len(users), len(test.userNames))
		} else {
			assert.Equal(t, len(users), len(test.userNames))
		}
	}
}

func TestGetEscalationPoliciesByTagIntegration(t *testing.T) {
	if os.Getenv("RUNNING_IN_JENKINS") == "true" {
		t.Skip("Skipping mPagerDuty TestGetEscalationPoliciesByTagIntegration in CI")
	}

	tests := []struct {
		tagID       string
		expectedErr bool
	}{
		{
			"P74RRGF", // GSOC tag
			false,
		},
		{
			"",
			true,
		},
		{
			"   ",
			true,
		},
		{
			"TEST",
			true,
		},
	}

	mPD, _ := mPagerDuty.GetMPagerDutyClient(authtoken)
	for _, test := range tests {
		response, err := mPD.GetEscalationPoliciesByTag(test.tagID)
		if test.expectedErr {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err)
			assert.Greater(t, response.EscalationPolicies, 0)
		}
	}
}

func TestUpdateEscalationPolicyIntegration(t *testing.T) {
	if os.Getenv("RUNNING_IN_JENKINS") == "true" {
		t.Skip("Skipping mPagerDuty TestGetOnCallUserIntegration in CI")
	}

	tests := []struct {
		id          string
		userID      string
		serviceID   string
		teamID      string
		escalation  []pagerduty.APIObject
		expectedErr bool
	}{
		{
			"P23N6LT", // INCY DEV Escalation Policy
			"PWKNFGT", // Timur Kalandarov
			"P39MEWZ", // GSOC General Service
			"P83EOFI", // GSOC Team
			[]pagerduty.APIObject{
				{
					ID:   "P273W1N", // Caleb Young
					Type: "user_reference",
				},
			},
			false,
		},
		{
			"DUMMYSTRING",
			"DUMMYSTRING",
			"DUMMYSTRING",
			"DUMMYSTRING",
			[]pagerduty.APIObject{
				{
					ID:   "DUMMYSTRING",
					Type: "user_reference",
				},
			},
			true,
		},
	}

	mPD, _ := mPagerDuty.GetMPagerDutyClient(authtoken)
	for _, test := range tests {
		escalationPolicy, err := mPD.UpdateEscalationPolicy(test.id, test.userID, test.serviceID, test.teamID, test.escalation)
		if test.expectedErr {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err)
			assert.Equal(t, escalationPolicy.ID, test.id)
		}
	}
}

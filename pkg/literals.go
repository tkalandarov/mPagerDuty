package mPagerDuty

import "os"

const limit uint = 100
const increase uint = 100

var pdScheduleTitlePrefix = os.Getenv("PD_SCHEDULEPREFIX")
var pgGSOCTeamID = []string{os.Getenv("PD_TEAMID")}

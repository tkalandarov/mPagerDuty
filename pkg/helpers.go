package mPagerDuty

import (
	"os"
	"unicode"

	"github.com/PagerDuty/go-pagerduty"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

type client struct {
	pdClient *pagerduty.Client
}

func newMPagerDutyClient(authtoken string) (IMPagerDuty, error) {
	if os.Getenv("RUNNING_IN_JENKINS") == "true" || os.Getenv("LOCAL_DEV_TESTING") == "true" {
		return &FakePDClient{}, nil
	}

	return &client{pdClient: pagerduty.NewClient(authtoken)}, nil
}

// Replaces non ASCII (accents, ąčęėįšųūž, etc...) characters with ASCII characters
func normalizeString(text string) (string, error) {
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	result, _, err := transform.String(t, text)
	return result, err
}

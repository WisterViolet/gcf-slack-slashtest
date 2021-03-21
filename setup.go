package slashtest

import (
	"os"
)

var (
	signingSecret string
)

func setup() {
	signingSecret = os.Getenv("SLACK_SIGNING_SECRET")
}

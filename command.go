package slashtest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type oldTimeStampError struct {
	s string
}

func (e *oldTimeStampError) Error() string {
	return e.s
}

const (
	version                     = "v0"
	slackRequestTimestampHeader = "X-Slack-Request-Timestamp"
	slackSignatureHeader        = "X-Slack-Signature"
)

func HelloBoard(w http.ResponseWriter, r *http.Request) {
	setup()

	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatalf("Coukdn't read request body: %v", err)
	}
	r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

	if r.Method != "POST" {
		http.Error(w, "Only POST requests are accepted", 405)
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Couldn't parse form", 400)
		log.Fatalf("ParseForm: %v", err)
	}

	r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	result, err := verifyRequest(r, signingSecret)
	if err != nil {
		log.Fatalf("verifyRequest: %v", err)
	}
	if !result {
		log.Fatalf("signatures did nor match.")
	}

	message := &Message{
		ResponseType: "in_channel",
		Text:         fmt.Sprintf("x-)/[Hello <@%s>]", r.Form["user_id"]),
		Attachments:  nil,
	}

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(message); err != nil {
		log.Fatalf("json.Marshal: %v", err)
	}
}

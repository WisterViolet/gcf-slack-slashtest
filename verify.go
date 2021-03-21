package slashtest

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func verifyRequest(r *http.Request, slackSigningSecret string) (bool, error) {
	timestamp := r.Header.Get(slackRequestTimestampHeader)
	slackSignature := r.Header.Get(slackSignatureHeader)

	t, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		return false, fmt.Errorf("strconv.ParseInt(%s): %v", timestamp, err)
	}

	if ageOk, age := checkTimestamp(t); !ageOk {
		return false, &oldTimeStampError{fmt.Sprintf("checkTimestamp(%v): %v %v", t, ageOk, age)}
	}

	if timestamp == "" || slackSignature == "" {
		return false, fmt.Errorf("ioutil.ReadAll")
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return false, fmt.Errorf("ioutil.ReadAll(%v): %v", r.Body, err)
	}

	r.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	baseString := fmt.Sprintf("%s:%s:%s", version, timestamp, body)

	signature := getSignature([]byte(baseString), []byte(slackSigningSecret))

	trimmed := strings.TrimPrefix(slackSignature, fmt.Sprintf("%s=", version))
	signatureInHeader, err := hex.DecodeString(trimmed)

	return hmac.Equal(signature, signatureInHeader), nil
}

func getSignature(base []byte, secret []byte) []byte {
	h := hmac.New(sha256.New, secret)
	h.Write(base)

	return h.Sum(nil)
}

func checkTimestamp(timestamp int64) (bool, time.Duration) {
	t := time.Since(time.Unix(timestamp, 0))

	return t.Minutes() <= 5, t
}

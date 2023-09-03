package couchbase

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"
	"unicode"

	"github.com/couchbase/gocb/v2"

	"github.com/trufflesecurity/trufflehog/v3/pkg/common"
	"github.com/trufflesecurity/trufflehog/v3/pkg/detectors"
	"github.com/trufflesecurity/trufflehog/v3/pkg/pb/detectorspb"
)

type Scanner struct{}

var _ detectors.Detector = (*Scanner)(nil)

var (

	// Make sure that your group is surrounded in boundary characters such as below to reduce false positives.
	connectionStringPat = regexp.MustCompile(`\bcb\.[a-z0-9]+\.cloud\.couchbase\.com\b`)
	usernamePat         = `?()/\+=\s\n`
	passwordPat         = `^<>;.*&|£\n\s`
	// passwordPat         = regexp.MustCompile(`(?i)(?:pass|pwd)(?:.|[\n\r]){0,15}(\b[^<>;.*&|£\n\s]{8,100}$)`)
	// passwordPat = regexp.MustCompile(`(?im)(?:pass|pwd)\S{0,40}?[:=\s]{1,3}[ '"=]{0,1}([^:?()/\+=\s\n]{4,40})\b`)
)

func meetsCouchbasePasswordRequirements(password []byte) ([]byte, bool) {
	var hasLower, hasUpper, hasNumber, hasSpecialChar bool
	passwordStr := string(password)
	for _, char := range passwordStr {
		switch {
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecialChar = true
		}

		if hasLower && hasUpper && hasNumber && hasSpecialChar {
			return password, true
		}
	}

	return nil, false
}

func (s Scanner) Keywords() [][]byte {
	return [][]byte{[]byte("couchbase://"), []byte("couchbases://")}
}

func (s Scanner) FromData(ctx context.Context, verify bool, data []byte) (results []detectors.Result, err error) {
	dataStr := string(data)

	connectionStringMatches := connectionStringPat.FindAllStringSubmatch(dataStr, -1)

	// prepend 'couchbases://' to the connection string as the connection
	// string format is couchbases://cb.stuff.cloud.couchbase.com but the
	// cb.stuff.cloud.couchbase.com may be separated from the couchbases:// in codebases.
	for i, connectionStringMatch := range connectionStringMatches {
		connectionStringMatches[i][0] = "couchbases://" + connectionStringMatch[0]
	}

	usernameRegexState := common.UsernameRegexCheck(usernamePat)
	usernameMatches := usernameRegexState.Matches(data)

	passwordRegexState := common.PasswordRegexCheck(passwordPat)
	passwordMatches := passwordRegexState.Matches(data)

	for _, connectionStringMatch := range connectionStringMatches {
		resConnectionStringMatch := strings.TrimSpace(connectionStringMatch[0])

		for _, resUsernameMatch := range usernameMatches {

			for _, resPasswordMatch := range passwordMatches {
				_, metPasswordRequirements := meetsCouchbasePasswordRequirements([]byte(resPasswordMatch))

				if !metPasswordRequirements {
					continue
				}

				s1 := detectors.Result{
					DetectorType: detectorspb.DetectorType_Couchbase,
					Raw:          []byte(fmt.Sprintf("%s:%s@%s", resUsernameMatch, resPasswordMatch, resConnectionStringMatch)),
				}

				if verify {

					options := gocb.ClusterOptions{
						Authenticator: gocb.PasswordAuthenticator{
							Username: resUsernameMatch,
							Password: resPasswordMatch,
						},
					}

					// Sets a pre-configured profile called "wan-development" to help avoid latency issues
					// when accessing Capella from a different Wide Area Network
					// or Availability Zone (e.g. your laptop).
					if err := options.ApplyProfile(gocb.ClusterConfigProfileWanDevelopment); err != nil {
						log.Fatal("apply profile err", err)
					}

					// Initialize the Connection
					cluster, err := gocb.Connect(resConnectionStringMatch, options)
					if err != nil {
						continue
					}

					// We'll ping the KV nodes in our cluster.
					pings, err := cluster.Ping(&gocb.PingOptions{
						Timeout: time.Second * 5,
					})

					if err != nil {
						continue
					}

					for _, ping := range pings.Services {
						for _, pingEndpoint := range ping {
							if pingEndpoint.State == gocb.PingStateOk {
								s1.Verified = true
								break
							} else {
								// This function will check false positives for common test words, but also it will make sure the key appears 'random' enough to be a real key.
								if detectors.IsKnownFalsePositive([]byte(resPasswordMatch), detectors.DefaultFalsePositives, true) {
									continue
								}
							}
						}
					}
				}

				results = append(results, s1)
			}
		}
	}
	return results, nil
}

func (s Scanner) Type() detectorspb.DetectorType {
	return detectorspb.DetectorType_Couchbase
}

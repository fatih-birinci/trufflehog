package dockerhub

import (
	"context"
	"fmt"
	regexp "github.com/wasilibs/go-re2"
	"io"
	"net/http"
	"strings"

	"github.com/trufflesecurity/trufflehog/v3/pkg/common"
	"github.com/trufflesecurity/trufflehog/v3/pkg/detectors"
	"github.com/trufflesecurity/trufflehog/v3/pkg/pb/detectorspb"
)

type Scanner struct{
	client *http.Client
}


// Ensure the Scanner satisfies the interfaces at compile time.
var _ detectors.Detector = (*Scanner)(nil)

var (
	defaultClient = common.SaneHttpClient()

	// Can use email or username for login.
	usernamePat = regexp.MustCompile(`(?im)(?:user|usr|-u|id)\S{0,40}?[:=\s]{1,3}[ '"=]?([a-zA-Z0-9]{4,40})\b`)
	emailPat    = regexp.MustCompile(common.EmailPattern)

	// Can use password or personal access token (PAT) for login, but this scanner will only check for PATs.
	accessTokenPat = regexp.MustCompile(detectors.PrefixRegex([]string{"docker"}) + `\b([0-9Aa-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12})\b`)
)

// Keywords are used for efficiently pre-filtering chunks.
// Use identifiers in the secret preferably, or the provider name.
func (s Scanner) Keywords() []string {
	return []string{"docker"}
}

// FromData will find and optionally verify Dockerhub secrets in a given set of bytes.
func (s Scanner) FromData(ctx context.Context, verify bool, data []byte) (results []detectors.Result, err error) {
	dataStr := string(data)

	emailMatches := emailPat.FindAllString(dataStr, -1)
	dataStr = emailPat.ReplaceAllString(dataStr, "")
	usernameMatches := usernamePat.FindAllStringSubmatch(dataStr, -1)

	accessTokenMatches := accessTokenPat.FindAllStringSubmatch(dataStr, -1)

	userMatches := emailMatches
	for _, usernameMatch := range usernameMatches {
		if len(usernameMatch) > 1 {
			userMatches = append(userMatches, usernameMatch[1])
		}
	}

	for _, resUserMatch := range userMatches {
		for _, resAccessTokenMatch := range accessTokenMatches {
			if len(resAccessTokenMatch) != 2 {
                        	continue
                	}
                	pat := strings.TrimSpace(resAccessTokenMatch[1])

			s1 := detectors.Result{
				DetectorType: detectorspb.DetectorType_Dockerhub,
				Raw:          []byte(fmt.Sprintf("%s: %s", resUserMatch, pat)),
			}

			if verify {
				client := s.client
                        	if client == nil {
                                	client = defaultClient
                        	}

				payload := strings.NewReader(fmt.Sprintf(`{"username": "%s", "password": "%s"}`, resUserMatch, pat))

				req, err := http.NewRequestWithContext(ctx, "GET", "https://hub.docker.com/v2/users/login", payload)
				if err != nil {
					continue
				}
				req.Header.Add("Content-Type", "application/json")
				res, err := client.Do(req)
				if err == nil {
					defer res.Body.Close()
					body, err := io.ReadAll(res.Body)
					if err != nil {
						continue
					}

					// Valid credentials can still return a 401 status code if 2FA is enabled
					if (res.StatusCode >= 200 && res.StatusCode < 300) || (res.StatusCode == 401 && strings.Contains(string(body), "login_2fa_token")) {
						s1.Verified = true
					} else if res.StatusCode == 401 {
                                        // The secret is determinately not verified (nothing to do)
                                	} else {
						s1.VerificationError = fmt.Errorf("unexpected HTTP response status %d", res.StatusCode)
					}
				} else {
        	                        s1.VerificationError = err
	                        }

			}

			results = append(results, s1)
		}
	}

	return results, nil
}

func (s Scanner) Type() detectorspb.DetectorType {
	return detectorspb.DetectorType_Dockerhub
}

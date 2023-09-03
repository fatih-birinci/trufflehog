package dockerhub

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"regexp"

	"github.com/trufflesecurity/trufflehog/v3/pkg/common"
	"github.com/trufflesecurity/trufflehog/v3/pkg/detectors"
	"github.com/trufflesecurity/trufflehog/v3/pkg/pb/detectorspb"
)

type Scanner struct{}

// Ensure the Scanner satisfies the interface at compile time.
var _ detectors.Detector = (*Scanner)(nil)

var (
	client = common.SaneHttpClient()

	usernamePat = regexp.MustCompile(`(?im)(?:user|usr|-u)\S{0,40}?[:=\s]{1,3}[ '"=]?([a-zA-Z0-9]{4,40})\b`)
	emailPat    = regexp.MustCompile(common.EmailPattern)

	// Can use password or personal access token (PAT) for login, but this scanner will only check for PATs.
	accessTokenPat = regexp.MustCompile(`\bdckr_pat_([a-zA-Z0-9_-]){27}\b`)
)

// Keywords are used for efficiently pre-filtering chunks.
func (s Scanner) Keywords() [][]byte {
	return [][]byte{[]byte("dckr_pat_")}
}

// FromData will find and optionally verify Dockerhub secrets in a given set of bytes.
func (s Scanner) FromData(ctx context.Context, verify bool, data []byte) (results []detectors.Result, err error) {

	emailMatches := emailPat.FindAll(data, -1)
	data = emailPat.ReplaceAll(data, []byte(""))

	usernameMatches := usernamePat.FindAllSubmatch(data, -1)
	accessTokenMatches := accessTokenPat.FindAll(data, -1)

	var userMatches [][]byte
	userMatches = append(userMatches, emailMatches...)
	for _, usernameMatch := range usernameMatches {
		if len(usernameMatch) > 1 {
			userMatches = append(userMatches, usernameMatch[1])
		}
	}

	for _, resUserMatch := range userMatches {
		for _, resAccessTokenMatch := range accessTokenMatches {

			s1 := detectors.Result{
				DetectorType: detectorspb.DetectorType_Dockerhub,
				Raw:          append(resUserMatch, append([]byte(": "), resAccessTokenMatch...)...),
			}

			if verify {
				payload := bytes.NewBuffer([]byte(fmt.Sprintf(`{"username": "%s", "password": "%s"}`, resUserMatch, resAccessTokenMatch)))

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
					if (res.StatusCode >= 200 && res.StatusCode < 300) || (res.StatusCode == 401 && bytes.Contains(body, []byte("login_2fa_token"))) {
						s1.Verified = true
					}
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

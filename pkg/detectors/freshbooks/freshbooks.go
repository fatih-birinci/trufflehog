package freshbooks

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"regexp"

	"io"

	"github.com/trufflesecurity/trufflehog/v3/pkg/common"
	"github.com/trufflesecurity/trufflehog/v3/pkg/detectors"
	"github.com/trufflesecurity/trufflehog/v3/pkg/pb/detectorspb"
)

type Scanner struct{}

// Ensure the Scanner satisfies the interface at compile time.
var _ detectors.Detector = (*Scanner)(nil)

var (
	client = common.SaneHttpClient()

	// Make sure that your group is surrounded in boundary characters such as below to reduce false positives.
	keyPat = regexp.MustCompile(detectors.PrefixRegex([]string{"freshbooks"}) + `\b([0-9a-z]{64})\b`)
	uriPat = regexp.MustCompile(detectors.PrefixRegex([]string{"freshbooks"}) + `\b(https://www.[0-9A-Za-z_-]{1,}.com)\b`)
)

// Keywords are used for efficiently pre-filtering chunks.
// Use identifiers in the secret preferably, or the provider name.
func (s Scanner) Keywords() [][]byte {
	return [][]byte{[]byte("freshbooks")}
}

// FromData will find and optionally verify Freshbooks secrets in a given set of bytes.
func (s Scanner) FromData(ctx context.Context, verify bool, data []byte) (results []detectors.Result, err error) {
	matches := keyPat.FindAllSubmatch(data, -1)
	uriMatches := uriPat.FindAllSubmatch(data, -1)

	for _, match := range matches {
		if len(match) != 2 {
			continue
		}
		resMatch := bytes.TrimSpace(match[1])
		for _, uriMatch := range uriMatches {
			if len(uriMatch) != 2 {
				continue
			}
			resURI := bytes.TrimSpace(uriMatch[1])

			s1 := detectors.Result{
				DetectorType: detectorspb.DetectorType_Freshbooks,
				Raw:          resMatch,
			}

			if verify {
				req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf(`https://auth.freshbooks.com/oauth/authorize?client_id=%s&redirect_uri=%s&response_type=code`, string(resMatch), string(resURI)), nil)
				if err != nil {
					continue
				}
				res, err := client.Do(req)
				if err == nil {
					defer res.Body.Close()
					bodyBytes, err := io.ReadAll(res.Body)
					if err != nil {
						continue
					}
					if res.StatusCode >= 200 && res.StatusCode < 300 && bytes.Contains(bodyBytes, []byte("Log In to FreshBooks")) {
						s1.Verified = true
					} else {
						if detectors.IsKnownFalsePositive(resMatch, detectors.DefaultFalsePositives, true) {
							continue
						}
					}
				}
			}

			results = append(results, s1)
		}
	}

	return results, nil
}

func (s Scanner) Type() detectorspb.DetectorType {
	return detectorspb.DetectorType_Freshbooks
}

package nexmoapikey

import (
	"bytes"
	"context"
	"net/http"
	"regexp"

	"github.com/trufflesecurity/trufflehog/v3/pkg/common"
	"github.com/trufflesecurity/trufflehog/v3/pkg/detectors"
	"github.com/trufflesecurity/trufflehog/v3/pkg/pb/detectorspb"
)

type Scanner struct{}

var _ detectors.Detector = (*Scanner)(nil)

var (
	client    = common.SaneHttpClient()
	keyPat    = regexp.MustCompile(detectors.PrefixRegex([]string{"nexmo"}) + `\b([A-Za-z0-9_-]{8})\b`)
	secretPat = regexp.MustCompile(detectors.PrefixRegex([]string{"nexmo"}) + `\b([A-Za-z0-9_-]{16})\b`)
)

func (s Scanner) Keywords() [][]byte {
	return [][]byte{[]byte("nexmo")}
}

func (s Scanner) FromData(ctx context.Context, verify bool, data []byte) (results []detectors.Result, err error) {
	matches := keyPat.FindAllSubmatch(data, -1)
	secretMatches := secretPat.FindAllSubmatch(data, -1)

	for _, match := range matches {
		if len(match) != 2 {
			continue
		}

		resMatch := bytes.TrimSpace(match[1])

		for _, secretMatch := range secretMatches {
			if len(secretMatch) != 2 {
				continue
			}

			resSecret := bytes.TrimSpace(secretMatch[1])

			s1 := detectors.Result{
				DetectorType: detectorspb.DetectorType_NexmoApiKey,
				Raw:          resMatch,
				RawV2:        append(resMatch, resSecret...),
			}

			if verify {
				req, err := http.NewRequestWithContext(ctx, "GET", "https://rest.nexmo.com/account/get-balance?api_key="+string(resMatch)+"&api_secret="+string(resSecret), nil)
				if err != nil {
					continue
				}

				res, err := client.Do(req)

				if err == nil {
					defer res.Body.Close()

					if res.StatusCode >= 200 && res.StatusCode < 300 {
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
	return detectorspb.DetectorType_NexmoApiKey
}

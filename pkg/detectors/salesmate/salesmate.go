package salesmate

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"regexp"

	"github.com/trufflesecurity/trufflehog/v3/pkg/common"
	"github.com/trufflesecurity/trufflehog/v3/pkg/detectors"
	"github.com/trufflesecurity/trufflehog/v3/pkg/pb/detectorspb"
)

type Scanner struct{}

var _ detectors.Detector = (*Scanner)(nil)

var (
	client = common.SaneHttpClient()

	keyPat    = regexp.MustCompile(detectors.PrefixRegex([]string{"salesmate"}) + `\b([0-9Aa-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12})\b`)
	domainPat = regexp.MustCompile(detectors.PrefixRegex([]string{"salesmate"}) + `\b([a-z0-9A-Z]{3,22})\b`)
)

func (s Scanner) Keywords() [][]byte {
	return [][]byte{[]byte("salesmate")}
}

func (s Scanner) FromData(ctx context.Context, verify bool, data []byte) (results []detectors.Result, err error) {

	matches := keyPat.FindAllSubmatch(data, -1)
	idmatches := domainPat.FindAllSubmatch(data, -1)

	for _, match := range matches {
		if len(match) != 2 {
			continue
		}
		resMatch := bytes.TrimSpace(match[1])
		for _, idmatch := range idmatches {
			if len(idmatch) != 2 {
				continue
			}
			resIdMatch := bytes.TrimSpace(idmatch[1])
			s1 := detectors.Result{
				DetectorType: detectorspb.DetectorType_Salesmate,
				Raw:          resMatch,
			}

			if verify {
				url := fmt.Sprintf("https://%s.salesmate.io/apis/v3/companies/1?trackingRecentSearch=true", string(resIdMatch))
				req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
				if err != nil {
					continue
				}

				req.Header.Add("Content-Type", "application/json")
				req.Header.Add("sessionToken", string(resMatch))
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
	return detectorspb.DetectorType_Salesmate
}

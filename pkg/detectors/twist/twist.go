package twist

import (
	"bytes"
	"context"
	"io"
	"mime/multipart"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/trufflesecurity/trufflehog/v3/pkg/common"
	"github.com/trufflesecurity/trufflehog/v3/pkg/detectors"
	"github.com/trufflesecurity/trufflehog/v3/pkg/pb/detectorspb"
)

type Scanner struct{}

// Ensure the Scanner satisfies the interface at compile time
var _ detectors.Detector = (*Scanner)(nil)

var (
	client = common.SaneHttpClient()

	//Make sure that your group is surrounded in boundry characters such as below to reduce false positives
	keyPat  = regexp.MustCompile(detectors.PrefixRegex([]string{"twist"}) + `\b([a-z0-9]{4,25}@[a-zA-Z0-9]{2,12}.[a-zA-Z0-9]{2,6})\b`)
	passPat = regexp.MustCompile(detectors.PrefixRegex([]string{"twist pass"}) + `\b([0-9A-Za-z\S]{7,20})\b`)
)

// Keywords are used for efficiently pre-filtering chunks.
// Use identifiers in the secret preferably, or the provider name.
func (s Scanner) Keywords() []string {
	return []string{"twist"}
}

// FromData will find and optionally verify Twist secrets in a given set of bytes.
func (s Scanner) FromData(ctx context.Context, verify bool, data []byte) (results []detectors.Result, err error) {
	dataStr := string(data)

	matches := keyPat.FindAllStringSubmatch(dataStr, -1)
	passMatches := passPat.FindAllStringSubmatch(dataStr, -1)

	for _, match := range matches {
		if len(match) != 2 {
			continue
		}
		resMatch := strings.TrimSpace(match[1])

		for _, passMatch := range passMatches {
			if len(passMatch) != 2 {
				continue
			}

			resPassMatch := strings.TrimSpace(passMatch[1])

			s1 := detectors.Result{
				DetectorType: detectorspb.DetectorType_Twist,
				Raw:          []byte(resMatch),
			}

			if verify {
				timeout := 10 * time.Second
				client.Timeout = timeout
				body := &bytes.Buffer{}
				writer := multipart.NewWriter(body)
				fw, err := writer.CreateFormField("email")
				if err != nil {
					continue
				}
				_, err = io.Copy(fw, strings.NewReader(resMatch))
				if err != nil {
					continue
				}
				fw, err = writer.CreateFormField("password")
				if err != nil {
					continue
				}
				_, err = io.Copy(fw, strings.NewReader(resPassMatch))
				if err != nil {
					continue
				}
				writer.Close()
				req, err := http.NewRequestWithContext(ctx, "POST", "https://api.twist.com/api/v3/users/login", bytes.NewReader(body.Bytes()))
				if err != nil {
					continue
				}
				req.Header.Add("Content-Type", writer.FormDataContentType())
				res, err := client.Do(req)
				if err == nil {
					defer res.Body.Close()
					if res.StatusCode >= 200 && res.StatusCode < 300 {
						s1.Verified = true
					} else {
						//This function will check false positives for common test words, but also it will make sure the key appears 'random' enough to be a real key
						if detectors.IsKnownFalsePositive(resPassMatch, detectors.DefaultFalsePositives, true) {
							continue
						}
					}
				}
			}
			results = append(results, s1)
		}
	}
	return detectors.CleanResults(results), nil
}

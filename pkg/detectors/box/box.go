package box

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	regexp "github.com/wasilibs/go-re2"

	"github.com/trufflesecurity/trufflehog/v3/pkg/common"
	"github.com/trufflesecurity/trufflehog/v3/pkg/detectors"
	"github.com/trufflesecurity/trufflehog/v3/pkg/pb/detectorspb"
)

type Scanner struct {
	client *http.Client
}

// Ensure the Scanner satisfies the interface at compile time.
var _ detectors.Detector = (*Scanner)(nil)

var (
	defaultClient = common.SaneHttpClient()
	// Make sure that your group is surrounded in boundary characters such as below to reduce false positives.
	keyPat = regexp.MustCompile(detectors.PrefixRegex([]string{"box", "access"}) + `\b([0-9a-zA-Z]{32})\b`)
)

// Keywords are used for efficiently pre-filtering chunks.
// Use identifiers in the secret preferably, or the provider name.
func (s Scanner) Keywords() []string {
	return []string{"box"}
}

// FromData will find and optionally verify Eraser secrets in a given set of bytes.
func (s Scanner) FromData(ctx context.Context, verify bool, data []byte) (results []detectors.Result, err error) {
	dataStr := string(data)

	uniqueMatches := make(map[string]struct{})
	for _, match := range keyPat.FindAllStringSubmatch(dataStr, -1) {
		uniqueMatches[match[1]] = struct{}{}
	}

	for match := range uniqueMatches {
		s1 := detectors.Result{
			DetectorType: detectorspb.DetectorType_Box,
			Raw:          []byte(match),
		}

		if verify {
			client := s.client
			if client == nil {
				client = defaultClient
			}

			isVerified, extraData, verificationErr := verifyMatch(ctx, client, match)
			s1.Verified = isVerified
			s1.ExtraData = extraData
			s1.SetVerificationError(verificationErr, match)
		}

		results = append(results, s1)
	}

	return
}

func verifyMatch(ctx context.Context, client *http.Client, token string) (bool, map[string]string, error) {
	url := "https://api.box.com/2.0/users/me"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return false, nil, nil
	}

	req.Header = http.Header{"Authorization": []string{"Bearer " + token}}
	req.Header.Add("content-type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return false, nil, err
	}
	defer func() {
		_, _ = io.Copy(io.Discard, res.Body)
		_ = res.Body.Close()
	}()

	switch res.StatusCode {
	case http.StatusOK:
		{
			var u user
			if err := json.NewDecoder(res.Body).Decode(&u); err != nil {
				return false, nil, err
			}
			return true, bakeExtraDataFromUser(u), nil
		}
	case http.StatusUnauthorized:
		// 401 access token not found
		// The secret is determinately not verified (nothing to do)
		return false, nil, nil
	default:
		return false, nil, fmt.Errorf("unexpected HTTP response status %d", res.StatusCode)
	}
}

func (s Scanner) Type() detectorspb.DetectorType {
	return detectorspb.DetectorType_Box
}

func bakeExtraDataFromUser(u user) map[string]string {
	return map[string]string{
		"user_id":           u.ID,
		"user_name":         u.Name,
		"username":          u.Login,
		"user_timezone":     u.Timezone,
		"user_space_amount": fmt.Sprintf("%d", u.SpaceAmount),
		"user_space_used":   fmt.Sprintf("%d", u.SpaceUsed),
		"user_status":       u.Status,
		"user_job_title":    u.JobTitle,
		"user_phone":        u.Phone,
		"user_address":      u.Address,
	}
}

// struct to represent a Box user.
type user struct {
	Type        string `json:"type"`
	ID          string `json:"id"`
	Name        string `json:"name"`
	Login       string `json:"login"`
	Timezone    string `json:"timezone"`
	SpaceAmount int64  `json:"space_amount"`
	SpaceUsed   int64  `json:"space_used"`
	Status      string `json:"status"`
	JobTitle    string `json:"job_title"`
	Phone       string `json:"phone"`
	Address     string `json:"address"`
}

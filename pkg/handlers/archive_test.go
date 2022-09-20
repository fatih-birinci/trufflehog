package handlers

import (
	"net/http"
	"regexp"
	"testing"

	diskbufferreader "github.com/bill-rich/disk-buffer-reader"
)

func TestArchiveHandler(t *testing.T) {
	tests := map[string]struct {
		archiveURL     string
		expectedChunks int
		matchString    string
	}{
		"gzip-single": {
			"https://raw.githubusercontent.com/bill-rich/bad-secrets/master/one-zip.gz",
			1,
			"AKIAYVP4CIPPH5TNP3SW",
		},
		"gzip-nested": {
			"https://raw.githubusercontent.com/bill-rich/bad-secrets/master/double-zip.gz",
			1,
			"AKIAYVP4CIPPH5TNP3SW",
		},
		"gzip-too-deep": {
			"https://raw.githubusercontent.com/bill-rich/bad-secrets/master/six-zip.gz",
			0,
			"",
		},
		"tar-single": {
			"https://raw.githubusercontent.com/bill-rich/bad-secrets/master/one.tar",
			1,
			"AKIAYVP4CIPPH5TNP3SW",
		},
		"tar-nested": {
			"https://raw.githubusercontent.com/bill-rich/bad-secrets/master/two.tar",
			1,
			"AKIAYVP4CIPPH5TNP3SW",
		},
		"tar-too-deep": {
			"https://raw.githubusercontent.com/bill-rich/bad-secrets/master/six.tar",
			0,
			"",
		},
		"targz-single": {
			"https://raw.githubusercontent.com/bill-rich/bad-secrets/master/tar-archive.tar.gz",
			1,
			"AKIAYVP4CIPPH5TNP3SW",
		},
		"gzip-large": {
			"https://raw.githubusercontent.com/bill-rich/bad-secrets/master/FifteenMB.gz",
			1543,
			"AKIAYVP4CIPPH5TNP3SW",
		},
		"zip-single": {
			"https://raw.githubusercontent.com/bill-rich/bad-secrets/master/aws-canary-creds.zip",
			1,
			"AKIAYVP4CIPPH5TNP3SW",
		},
	}

	for name, testCase := range tests {
		resp, err := http.Get(testCase.archiveURL)
		if err != nil || resp.StatusCode != http.StatusOK {
			t.Error(err)
		}
		defer resp.Body.Close()

		archive := Archive{}
		archive.New()

		newReader, err := diskbufferreader.New(resp.Body)
		if err != nil {
			t.Errorf("error creating reusable reader: %s", err)
		}
		archiveChan := archive.FromFile(newReader)

		count := 0
		re := regexp.MustCompile(testCase.matchString)
		matched := false
		for chunk := range archiveChan {
			count++
			if re.Match(chunk) {
				matched = true
			}
		}
		if !matched && len(testCase.matchString) > 0 {
			t.Errorf("%s: Expected string not found in archive.", name)
		}
		if count != testCase.expectedChunks {
			t.Errorf("%s: Unexpected number of chunks. Got %d, expected: %d", name, count, testCase.expectedChunks)
		}
	}
}

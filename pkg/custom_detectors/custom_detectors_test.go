package custom_detectors

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/trufflesecurity/trufflehog/v3/pkg/pb/custom_detectorspb"
	"github.com/trufflesecurity/trufflehog/v3/pkg/protoyaml"
)

func TestCustomRegexTemplateParsing(t *testing.T) {
	testCustomRegexTemplateYaml := `name: Internal bi tool
keywords:
- secret_v1_
- pat_v2_
regex:
  id_pat_example: ([a-zA-Z0-9]{32})
  secret_pat_example: ([a-zA-Z0-9]{32})
verify:
- endpoint: http://localhost:8000/{id_pat_example}
  unsafe: true
  headers:
  - 'Authorization: Bearer {secret_pat_example.0}'
  successRanges:
  - 200-250
  - '288'`

	var got custom_detectorspb.CustomRegex
	assert.NoError(t, protoyaml.UnmarshalStrict([]byte(testCustomRegexTemplateYaml), &got))
	assert.Equal(t, "Internal bi tool", got.Name)
	assert.Equal(t, []string{"secret_v1_", "pat_v2_"}, got.Keywords)
	assert.Equal(t, map[string]string{
		"id_pat_example":     "([a-zA-Z0-9]{32})",
		"secret_pat_example": "([a-zA-Z0-9]{32})",
	}, got.Regex)
	assert.Equal(t, 1, len(got.Verify))
	assert.Equal(t, "http://localhost:8000/{id_pat_example}", got.Verify[0].Endpoint)
	assert.Equal(t, true, got.Verify[0].Unsafe)
	assert.Equal(t, []string{"Authorization: Bearer {secret_pat_example.0}"}, got.Verify[0].Headers)
	assert.Equal(t, []string{"200-250", "288"}, got.Verify[0].SuccessRanges)
}

func TestCustomRegexWebhookParsing(t *testing.T) {
	testCustomRegexWebhookYaml := `name: Internal bi tool
keywords:
- secret_v1_
- pat_v2_
regex:
  id_pat_example: ([a-zA-Z0-9]{32})
  secret_pat_example: ([a-zA-Z0-9]{32})
verify:
- endpoint: http://localhost:8000/
  unsafe: true
  headers:
  - 'Authorization: Bearer token'`

	var got custom_detectorspb.CustomRegex
	assert.NoError(t, protoyaml.UnmarshalStrict([]byte(testCustomRegexWebhookYaml), &got))
	assert.Equal(t, "Internal bi tool", got.Name)
	assert.Equal(t, []string{"secret_v1_", "pat_v2_"}, got.Keywords)
	assert.Equal(t, map[string]string{
		"id_pat_example":     "([a-zA-Z0-9]{32})",
		"secret_pat_example": "([a-zA-Z0-9]{32})",
	}, got.Regex)
	assert.Equal(t, 1, len(got.Verify))
	assert.Equal(t, "http://localhost:8000/", got.Verify[0].Endpoint)
	assert.Equal(t, true, got.Verify[0].Unsafe)
	assert.Equal(t, []string{"Authorization: Bearer token"}, got.Verify[0].Headers)
}

// TestCustomDetectorsParsing tests the full `detectors` configuration.
func TestCustomDetectorsParsing(t *testing.T) {
	// TODO: Support both template and webhook.
	testYamlConfig := `detectors:
- name: Internal bi tool
  keywords:
  - secret_v1_
  - pat_v2_
  regex:
    id_pat_example: ([a-zA-Z0-9]{32})
    secret_pat_example: ([a-zA-Z0-9]{32})
  verify:
  - endpoint: http://localhost:8000/
    unsafe: true
    headers:
    - 'Authorization: Bearer token'`

	var messages custom_detectorspb.CustomDetectors
	assert.NoError(t, protoyaml.UnmarshalStrict([]byte(testYamlConfig), &messages))
	assert.Equal(t, 1, len(messages.Detectors))

	got := messages.Detectors[0]
	assert.Equal(t, "Internal bi tool", got.Name)
	assert.Equal(t, []string{"secret_v1_", "pat_v2_"}, got.Keywords)
	assert.Equal(t, map[string]string{
		"id_pat_example":     "([a-zA-Z0-9]{32})",
		"secret_pat_example": "([a-zA-Z0-9]{32})",
	}, got.Regex)
	assert.Equal(t, 1, len(got.Verify))
	assert.Equal(t, "http://localhost:8000/", got.Verify[0].Endpoint)
	assert.Equal(t, true, got.Verify[0].Unsafe)
	assert.Equal(t, []string{"Authorization: Bearer token"}, got.Verify[0].Headers)
}

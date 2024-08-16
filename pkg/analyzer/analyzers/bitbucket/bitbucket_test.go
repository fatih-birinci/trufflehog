package bitbucket

import (
	"encoding/json"
	"sort"
	"testing"
	"time"

	"github.com/trufflesecurity/trufflehog/v3/pkg/analyzer/analyzers"
	"github.com/trufflesecurity/trufflehog/v3/pkg/common"
	"github.com/trufflesecurity/trufflehog/v3/pkg/context"
)

func TestAnalyzer_Analyze(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	testSecrets, err := common.GetSecret(ctx, "trufflehog-testing", "detectors5")
	if err != nil {
		t.Fatalf("could not get test secrets from GCP: %s", err)
	}

	tests := []struct {
		name    string
		sid     string
		key     string
		want    string // JSON string
		wantErr bool
	}{
		{
			name:    "valid Bitbucket key",
			key:     testSecrets.MustGetField("BITBUCKET_ANALYZE_TOKEN"),
			want:    `{"AnalyzerType":3,"Bindings":[{"Resource":{"Name":"Repository","FullyQualifiedName":"bitbucket.com/credential/repo_access_token","Type":"repo_access_token","Metadata":{"type":"Repository Access Token (Can access 1 repository)"},"Parent":null},"Permission":{"Value":"pullrequest","Parent":null}},{"Resource":{"Name":"Repository","FullyQualifiedName":"bitbucket.com/credential/repo_access_token","Type":"repo_access_token","Metadata":{"type":"Repository Access Token (Can access 1 repository)"},"Parent":null},"Permission":{"Value":"webhook","Parent":null}},{"Resource":{"Name":"Repository","FullyQualifiedName":"bitbucket.com/credential/repo_access_token","Type":"repo_access_token","Metadata":{"type":"Repository Access Token (Can access 1 repository)"},"Parent":null},"Permission":{"Value":"pipeline","Parent":null}},{"Resource":{"Name":"Repository","FullyQualifiedName":"bitbucket.com/credential/repo_access_token","Type":"repo_access_token","Metadata":{"type":"Repository Access Token (Can access 1 repository)"},"Parent":null},"Permission":{"Value":"runner","Parent":null}}],"UnboundedResources":[{"Name":"bitbucket.com/repository/basit-trufflesec/repo1","FullyQualifiedName":"{8961ef70-000c-47ca-9348-5f9ecee875d6}","Type":"repository","Metadata":{"isPrivate":true,"owner":"basit-trufflesec","owner_id":"{521b49b6-7709-484a-8aa8-ecc3a6da08eb}","role":"admin"},"Parent":{"Name":"repo-analyzer","FullyQualifiedName":"bitbucket.com/project/{8a693e10-087f-41fc-ba67-2d1414ab1c86}","Type":"project","Metadata":null,"Parent":{"Name":"basit-trufflesec","FullyQualifiedName":"bitbucket.com/workspace/{521b49b6-7709-484a-8aa8-ecc3a6da08eb}","Type":"workspace","Metadata":null,"Parent":null}}}],"Metadata":null}`,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := Analyzer{}
			got, err := a.Analyze(ctx, map[string]string{"key": tt.key, "sid": tt.sid})
			if (err != nil) != tt.wantErr {
				t.Errorf("Analyzer.Analyze() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// bindings need to be in the same order to be comparable
			sortBindings(got.Bindings)

			// Marshal the actual result to JSON
			gotJSON, err := json.Marshal(got)
			if err != nil {
				t.Fatalf("could not marshal got to JSON: %s", err)
			}

			// Parse the expected JSON string
			var wantObj analyzers.AnalyzerResult
			if err := json.Unmarshal([]byte(tt.want), &wantObj); err != nil {
				t.Fatalf("could not unmarshal want JSON string: %s", err)
			}

			// bindings need to be in the same order to be comparable
			sortBindings(wantObj.Bindings)

			// Marshal the expected result to JSON (to normalize)
			wantJSON, err := json.Marshal(wantObj)
			if err != nil {
				t.Fatalf("could not marshal want to JSON: %s", err)
			}

			// Compare the JSON strings
			if string(gotJSON) != string(wantJSON) {
				// Pretty-print both JSON strings for easier comparison
				var gotIndented []byte
				gotIndented, err = json.MarshalIndent(got, "", " ")
				if err != nil {
					t.Fatalf("could not marshal got to indented JSON: %s", err)
				}
				t.Errorf("Analyzer.Analyze() = \n%s", gotIndented)
			}
		})
	}
}

// Helper function to sort bindings
func sortBindings(bindings []analyzers.Binding) {
	sort.SliceStable(bindings, func(i, j int) bool {
		if bindings[i].Resource.Name == bindings[j].Resource.Name {
			return bindings[i].Permission.Value < bindings[j].Permission.Value
		}
		return bindings[i].Resource.Name < bindings[j].Resource.Name
	})
}

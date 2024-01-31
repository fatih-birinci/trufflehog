package decoders

import (
	"testing"

	"github.com/kylelemons/godebug/pretty"

	"github.com/trufflesecurity/trufflehog/v3/pkg/sources"
)

func TestUnicodeEscape_FromChunk(t *testing.T) {
	tests := []struct {
		name    string
		chunk   *sources.Chunk
		want    *sources.Chunk
		wantErr bool
	}{
		{
			name: "all escaped",
			chunk: &sources.Chunk{
				Data: []byte("\\u0074\\u006f\\u006b\\u0065\\u006e\\u003a\\u0020\\u0022\\u0067\\u0068\\u0070\\u005f\\u0049\\u0077\\u0064\\u004d\\u0078\\u0039\\u0057\\u0046\\u0057\\u0052\\u0052\\u0066\\u004d\\u0068\\u0054\\u0059\\u0069\\u0061\\u0056\\u006a\\u005a\\u0037\\u0038\\u004a\\u0066\\u0075\\u0061\\u006d\\u0076\\u006e\\u0030\\u0059\\u0057\\u0052\\u004d\\u0030\\u0022"),
			},
			want: &sources.Chunk{
				Data: []byte("token: \"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\""),
			},
		},
		{
			name: "mixed content",
			chunk: &sources.Chunk{
				Data: []byte("npm config set @trufflesec:registry=https://npm.pkg.github.com\nnpm config set //npm.pkg.github.com:_authToken=$'\\u0067hp_9ovSHEBCq0drG42yjoam76iNybtqLN25CgSf'"),
			},
			want: &sources.Chunk{
				Data: []byte("npm config set @trufflesec:registry=https://npm.pkg.github.com\nnpm config set //npm.pkg.github.com:_authToken=$'ghp_9ovSHEBCq0drG42yjoam76iNybtqLN25CgSf'"),
			},
		},
		{
			name: "multiple slashes",
			chunk: &sources.Chunk{
				Data: []byte(`SameValue("hello","\\u0068el\\u006co");          // true`),
			},
			want: &sources.Chunk{
				Data: []byte(`SameValue("hello","hello");          // true`),
			},
		},
		{
			name: "no escaped",
			chunk: &sources.Chunk{
				Data: []byte(`-//npm.fontawesome.com/:_authToken=12345678-2323-1111-1111-12345670B312
+//npm.fontawesome.com/:_authToken=REMOVED_TOKEN`),
			},
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &EscapedUnicode{}
			got := d.FromChunk(tt.chunk)
			if tt.want != nil {
				if got == nil {
					t.Fatal("got nil, did not want nil")
				}
				if diff := pretty.Compare(string(tt.want.Data), string(got.Data)); diff != "" {
					t.Errorf("UnicodeEscape.FromChunk() %s diff: (-want +got)\n%s", tt.name, diff)
				}
			} else {
				if got != nil {
					t.Error("Expected nil chunk")
				}
			}
		})
	}
}

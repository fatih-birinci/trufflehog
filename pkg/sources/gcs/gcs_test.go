package gcs

import (
	"fmt"
	"io"
	"sort"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/trufflesecurity/trufflehog/v3/pkg/context"
	"github.com/trufflesecurity/trufflehog/v3/pkg/pb/source_metadatapb"
	"github.com/trufflesecurity/trufflehog/v3/pkg/pb/sourcespb"
	"github.com/trufflesecurity/trufflehog/v3/pkg/sources"
)

func createTestSource(src *sourcespb.GCS) (*Source, *anypb.Any) {
	s := &Source{}
	conn, err := anypb.New(src)
	if err != nil {
		panic(err)
	}
	return s, conn
}

func TestSourceInit(t *testing.T) {
	source, conn := createTestSource(&sourcespb.GCS{
		ProjectId: testProjectID,
		IncludeBuckets: []string{
			publicBucket,
		},
		ExcludeBuckets: []string{
			perfTestBucketGlob,
		},
		ExcludeObjects: []string{
			"object1",
		},
		Credential: &sourcespb.GCS_Unauthenticated{},
	})

	err := source.Init(context.Background(), "test", 1, 1, true, conn, 8)
	assert.Nil(t, err)
	assert.NotNil(t, source.gcsManager)
}

func TestConfigureGCSManager(t *testing.T) {
	testCases := []struct {
		name    string
		conn    *sourcespb.GCS
		want    *gcsManager
		wantErr bool
	}{
		{
			name:    "nil conn",
			wantErr: true,
		},
		{
			name: "valid conn, bare config",
			conn: &sourcespb.GCS{
				ProjectId:  testProjectID,
				Credential: &sourcespb.GCS_Adc{},
				ExcludeBuckets: []string{
					perfTestBucketGlob,
				},
			},
			want: &gcsManager{
				projectID:      testProjectID,
				excludeBuckets: map[string]struct{}{perfTestBucketGlob: {}},
			},
		},
		{
			name: "valid conn, include and exclude buckets",
			conn: &sourcespb.GCS{
				ProjectId:  testProjectID,
				Credential: &sourcespb.GCS_Adc{},
				IncludeBuckets: []string{
					"bucket1",
				},
				ExcludeBuckets: []string{
					perfTestBucketGlob,
				},
			},
			want: &gcsManager{
				projectID:      testProjectID,
				includeBuckets: map[string]struct{}{"bucket1": {}},
			},
		},
		{
			name: "valid conn, include and exclude objects",
			conn: &sourcespb.GCS{
				ProjectId:  testProjectID,
				Credential: &sourcespb.GCS_Adc{},
				IncludeObjects: []string{
					"object1",
				},
				ExcludeObjects: []string{
					"object2",
				},
				ExcludeBuckets: []string{
					perfTestBucketGlob,
				},
			},
			want: &gcsManager{
				projectID:      testProjectID,
				includeObjects: map[string]struct{}{"object1": {}},
				excludeBuckets: map[string]struct{}{perfTestBucketGlob: {}},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			got, err := configureGCSManager(ctx, tc.conn, 8)
			if err != nil && !tc.wantErr {
				t.Errorf("source.configureGCSManager() error = %v", err)
			}

			if !tc.wantErr {
				if diff := cmp.Diff(tc.want, got,
					cmp.AllowUnexported(gcsManager{}),
					cmpopts.IgnoreFields(gcsManager{}, "client", "workerPool", "concurrency", "buckets", "maxObjectSize", "attr"),
				); diff != "" {
					t.Errorf("source.Init() diff: (-want +got)\n%s", diff)
				}
			}
		})
	}
}

type mockObjectManager struct {
	wantErr bool
}

func (m *mockObjectManager) attributes(_ context.Context) (*attributes, error) {
	if m.wantErr {
		return nil, fmt.Errorf("some error")
	}

	return &attributes{
		numObjects:    5,
		numBuckets:    1,
		bucketObjects: map[string]uint64{testBucket: 5},
	}, nil
}

type mockReader struct {
	offset int
	data   []byte
}

func (m *mockReader) Read(p []byte) (n int, err error) {
	if m.offset >= len(m.data) {
		return 0, io.EOF
	}

	n = copy(p, m.data[m.offset:])
	m.offset += n
	return
}

func (m *mockObjectManager) listObjects(context.Context) (chan io.Reader, error) {
	if m.wantErr {
		return nil, fmt.Errorf("some error")
	}

	ch := make(chan io.Reader)
	go func() {
		defer close(ch)
		// Add 5 objects to the channel.
		for i := 0; i < 5; i++ {
			ch <- createTestObject(i)
		}
	}()

	return ch, nil
}

func createTestObject(id int) object {
	return object{
		name:        fmt.Sprintf("object%d", id),
		bucket:      testBucket,
		contentType: "plain/text",
		owner:       "testman@test.com",
		link:        fmt.Sprintf("https://storage.googleapis.com/%s/%s", testBucket, fmt.Sprintf("object%d", id)),
		acl:         []string{"authenticatedUsers"},
		size:        42,
		Reader:      &mockReader{data: []byte(fmt.Sprintf("hello world %d", id))},
		createdAt:   time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
	}
}

func createTestSourceChunk(id int) *sources.Chunk {
	return &sources.Chunk{
		SourceName: "test",
		SourceType: sourcespb.SourceType_SOURCE_TYPE_GCS,
		SourceID:   0,
		Verify:     true,
		Data:       []byte(fmt.Sprintf("hello world %d", id)),
		SourceMetadata: &source_metadatapb.MetaData{
			Data: &source_metadatapb.MetaData_Gcs{
				Gcs: &source_metadatapb.GCS{
					Filename:    fmt.Sprintf("object%d", id),
					Bucket:      testBucket,
					ContentType: "plain/text",
					Email:       "testman@test.com",
					Link:        fmt.Sprintf("https://storage.googleapis.com/%s/%s", testBucket, fmt.Sprintf("object%d", id)),
					Acls:        []string{"authenticatedUsers"},
					CreatedAt:   "1577836800",
				},
			},
		},
	}
}

func TestSourceChunks_ListObjects(t *testing.T) {
	ctx := context.Background()
	chunksCh := make(chan *sources.Chunk, 1)

	source := &Source{
		gcsManager: &mockObjectManager{},
		chunksCh:   chunksCh,
	}

	err := source.enumerate(ctx)
	assert.Nil(t, err)

	go func() {
		defer close(chunksCh)
		err := source.Chunks(ctx, chunksCh)
		assert.Nil(t, err)
	}()

	want := make([]*sources.Chunk, 0, 5)
	for i := 0; i < 5; i++ {
		want = append(want, createTestSourceChunk(i))
	}

	count := 0
	got := make([]*sources.Chunk, 0, 5)
	for ch := range chunksCh {
		got = append(got, ch)
		count++
	}

	// Ensure we get 5 objects back.
	assert.Equal(t, 5, count)

	// Sort the results to ensure deterministic ordering.
	sort.Slice(want, func(i, j int) bool {
		return want[i].SourceMetadata.GetGcs().Filename < want[j].SourceMetadata.GetGcs().Filename
	})
	sort.Slice(got, func(i, j int) bool {
		return got[i].SourceMetadata.GetGcs().Filename < got[j].SourceMetadata.GetGcs().Filename
	})

	for i, c := range got {
		assert.Equal(t, want[i].SourceMetadata.GetGcs().Filename, c.SourceMetadata.GetGcs().Filename)
		assert.Equal(t, want[i].Data, c.Data)
		assert.Equal(t, want[i].SourceMetadata.GetGcs().CreatedAt, c.SourceMetadata.GetGcs().CreatedAt)
	}

}

func TestSourceInit_Enumerate(t *testing.T) {
	ctx := context.Background()
	source := &Source{gcsManager: &mockObjectManager{}}

	err := source.enumerate(ctx)
	assert.Nil(t, err)

	// Ensure the attributes are set.
	assert.Equal(t, uint64(5), source.stats.numObjects)
	assert.Equal(t, uint32(1), source.stats.numBuckets)
	assert.Equal(t, uint64(5), source.stats.bucketObjects[testBucket])
}

func TestSourceChunks_ListObjects_Error(t *testing.T) {
	ctx := context.Background()
	source := &Source{gcsManager: &mockObjectManager{wantErr: true}}

	chunksCh := make(chan *sources.Chunk, 1)

	defer close(chunksCh)
	err := source.Chunks(ctx, chunksCh)
	assert.True(t, err != nil)
}

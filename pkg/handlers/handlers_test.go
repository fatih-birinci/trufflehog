package handlers

import (
	"archive/tar"
	"encoding/binary"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	diskbufferreader "github.com/trufflesecurity/disk-buffer-reader"
	"golang.org/x/exp/rand"

	logContext "github.com/trufflesecurity/trufflehog/v3/pkg/context"
	"github.com/trufflesecurity/trufflehog/v3/pkg/sources"
)

func TestHandleFileCancelledContext(t *testing.T) {
	reporter := sources.ChanReporter{Ch: make(chan *sources.Chunk, 2)}

	canceledCtx, cancel := logContext.WithCancel(logContext.Background())
	cancel()
	reader, err := diskbufferreader.New(strings.NewReader("file"))
	assert.NoError(t, err)
	assert.Error(t, HandleFile(canceledCtx, reader, &sources.Chunk{}, reporter))
}

func TestHandleFile(t *testing.T) {
	reporter := sources.ChanReporter{Ch: make(chan *sources.Chunk, 2)}

	// Only one chunk is sent on the channel.
	// TODO: Embed a zip without making an HTTP request.
	resp, err := http.Get("https://raw.githubusercontent.com/bill-rich/bad-secrets/master/aws-canary-creds.zip")
	assert.NoError(t, err)
	defer resp.Body.Close()

	reader, err := diskbufferreader.New(resp.Body)
	assert.NoError(t, err)

	assert.Equal(t, 0, len(reporter.Ch))
	assert.NoError(t, HandleFile(logContext.Background(), reader, &sources.Chunk{}, reporter))
	assert.Equal(t, 1, len(reporter.Ch))
}

func BenchmarkHandleFile(b *testing.B) {
	file, err := os.Open("testdata/test.tgz")
	assert.Nil(b, err)
	defer file.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sourceChan := make(chan *sources.Chunk, 1)
		reader, err := diskbufferreader.New(file)
		assert.NoError(b, err)

		b.StartTimer()

		go func() {
			defer close(sourceChan)
			err := HandleFile(logContext.Background(), reader, &sources.Chunk{}, sources.ChanReporter{Ch: sourceChan})
			assert.NoError(b, err)
		}()

		for range sourceChan {
		}

		b.StopTimer()
	}
}

func TestSkipArchive(t *testing.T) {
	file, err := os.Open("testdata/test.tgz")
	assert.Nil(t, err)
	defer file.Close()

	reader, err := diskbufferreader.New(file)
	assert.NoError(t, err)

	ctx := logContext.Background()

	chunkCh := make(chan *sources.Chunk)
	go func() {
		defer close(chunkCh)
		err := HandleFile(ctx, reader, &sources.Chunk{}, sources.ChanReporter{Ch: chunkCh}, WithSkipArchives(true))
		assert.NoError(t, err)
	}()

	wantCount := 0
	count := 0
	for range chunkCh {
		count++
	}
	assert.Equal(t, wantCount, count)
}

func TestHandleNestedArchives(t *testing.T) {
	file, err := os.Open("testdata/nested-dirs.zip")
	assert.Nil(t, err)
	defer file.Close()

	reader, err := diskbufferreader.New(file)
	assert.NoError(t, err)

	ctx := logContext.Background()

	chunkCh := make(chan *sources.Chunk)
	go func() {
		defer close(chunkCh)
		err := HandleFile(ctx, reader, &sources.Chunk{}, sources.ChanReporter{Ch: chunkCh})
		assert.NoError(t, err)
	}()

	wantCount := 8
	count := 0
	for range chunkCh {
		count++
	}
	assert.Equal(t, wantCount, count)
}

func TestHandleCompressedZip(t *testing.T) {
	file, err := os.Open("testdata/example.zip.gz")
	assert.Nil(t, err)
	defer file.Close()

	reader, err := diskbufferreader.New(file)
	assert.NoError(t, err)

	ctx := logContext.Background()

	chunkCh := make(chan *sources.Chunk)
	go func() {
		defer close(chunkCh)
		err := HandleFile(ctx, reader, &sources.Chunk{}, sources.ChanReporter{Ch: chunkCh})
		assert.NoError(t, err)
	}()

	wantCount := 2
	count := 0
	for range chunkCh {
		count++
	}
	assert.Equal(t, wantCount, count)
}

func TestHandleNestedCompressedArchive(t *testing.T) {
	file, err := os.Open("testdata/nested-compressed-archive.tar.gz")
	assert.Nil(t, err)
	defer file.Close()

	reader, err := diskbufferreader.New(file)
	assert.NoError(t, err)

	ctx := logContext.Background()

	chunkCh := make(chan *sources.Chunk)
	go func() {
		defer close(chunkCh)
		err := HandleFile(ctx, reader, &sources.Chunk{}, sources.ChanReporter{Ch: chunkCh})
		assert.NoError(t, err)
	}()

	wantCount := 4
	count := 0
	for range chunkCh {
		count++
	}
	assert.Equal(t, wantCount, count)
}

func TestExtractTarContent(t *testing.T) {
	file, err := os.Open("testdata/test.tgz")
	assert.Nil(t, err)
	defer file.Close()

	reader, err := diskbufferreader.New(file)
	assert.NoError(t, err)

	ctx := logContext.Background()

	chunkCh := make(chan *sources.Chunk)
	go func() {
		defer close(chunkCh)
		err := HandleFile(ctx, reader, &sources.Chunk{}, sources.ChanReporter{Ch: chunkCh})
		assert.NoError(t, err)
	}()

	wantCount := 4
	count := 0
	for range chunkCh {
		count++
	}
	assert.Equal(t, wantCount, count)
}

func TestNestedDirArchive(t *testing.T) {
	file, err := os.Open("testdata/dir-archive.zip")
	assert.Nil(t, err)
	defer file.Close()

	reader, err := diskbufferreader.New(file)
	assert.NoError(t, err)

	ctx, cancel := logContext.WithTimeout(logContext.Background(), 5*time.Second)
	defer cancel()
	sourceChan := make(chan *sources.Chunk, 1)

	go func() {
		defer close(sourceChan)
		err := HandleFile(ctx, reader, &sources.Chunk{}, sources.ChanReporter{Ch: sourceChan})
		assert.NoError(t, err)
	}()

	count := 0
	want := 4
	for range sourceChan {
		count++
	}
	assert.Equal(t, want, count)
}

func TestHandleFileRPM(t *testing.T) {
	wantChunkCount := 179
	reporter := sources.ChanReporter{Ch: make(chan *sources.Chunk, wantChunkCount)}

	file, err := os.Open("testdata/test.rpm")
	assert.Nil(t, err)
	defer file.Close()

	reader, err := diskbufferreader.New(file)
	assert.NoError(t, err)

	assert.Equal(t, 0, len(reporter.Ch))
	assert.NoError(t, HandleFile(logContext.Background(), reader, &sources.Chunk{}, reporter))
	assert.Equal(t, wantChunkCount, len(reporter.Ch))
}

func TestHandleFileAR(t *testing.T) {
	wantChunkCount := 102
	reporter := sources.ChanReporter{Ch: make(chan *sources.Chunk, wantChunkCount)}

	file, err := os.Open("testdata/test.deb")
	assert.Nil(t, err)
	defer file.Close()

	reader, err := diskbufferreader.New(file)
	assert.NoError(t, err)

	assert.Equal(t, 0, len(reporter.Ch))
	assert.NoError(t, HandleFile(logContext.Background(), reader, &sources.Chunk{}, reporter))
	assert.Equal(t, wantChunkCount, len(reporter.Ch))
}

func TestHandleFileSkipNonTextFiles(t *testing.T) {
	filename := createBinaryArchive(t)
	defer os.Remove(filename)

	file, err := os.Open(filename)
	assert.NoError(t, err)

	reader, err := diskbufferreader.New(file)
	assert.NoError(t, err)

	ctx, cancel := logContext.WithTimeout(logContext.Background(), 5*time.Second)
	defer cancel()
	sourceChan := make(chan *sources.Chunk, 1)

	go func() {
		defer close(sourceChan)
		err = HandleFile(ctx, reader, &sources.Chunk{}, sources.ChanReporter{Ch: sourceChan})
		assert.NoError(t, err)
	}()

	count := 0
	for range sourceChan {
		count++
	}
	// The binary archive should not be scanned.
	assert.Equal(t, 0, count)
}

func createBinaryArchive(t *testing.T) string {
	t.Helper()

	f, err := os.CreateTemp("", "testbinary")
	assert.NoError(t, err)
	defer os.Remove(f.Name())

	r := rand.New(rand.NewSource(uint64(time.Now().UnixNano())))

	randomBytes := make([]byte, 1024)
	_, err = r.Read(randomBytes)
	assert.NoError(t, err)

	_, err = f.Write(randomBytes)
	assert.NoError(t, err)

	// Create and write some structured binary data (e.g., integers, floats)
	for i := 0; i < 10; i++ {
		err = binary.Write(f, binary.LittleEndian, int32(rand.Intn(1000)))
		assert.NoError(t, err)
		err = binary.Write(f, binary.LittleEndian, rand.Float64())
		assert.NoError(t, err)
	}

	pngFile, err := os.CreateTemp("", "example.bin")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(pngFile.Name())

	// Write the binary content to the .bin file
	fileContent, err := os.ReadFile(f.Name())
	assert.NoError(t, err)

	_, err = pngFile.Write(fileContent)
	assert.NoError(t, err)

	tarFile, err := os.Create("example.tar")
	if err != nil {
		t.Fatal(err)
	}
	defer tarFile.Close()

	// Create a new tar archive.
	tarWriter := tar.NewWriter(tarFile)
	defer tarWriter.Close()

	fileInfo, err := pngFile.Stat()
	assert.NoError(t, err)

	header, err := tar.FileInfoHeader(fileInfo, "")
	assert.NoError(t, err)

	header.Name = "example.png"
	err = tarWriter.WriteHeader(header)
	assert.NoError(t, err)

	_, err = pngFile.Seek(0, 0)
	assert.NoError(t, err)

	_, err = io.Copy(tarWriter, pngFile)
	assert.NoError(t, err)

	return tarFile.Name()
}

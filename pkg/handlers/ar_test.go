package handlers

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/trufflesecurity/trufflehog/v3/pkg/context"
)

func TestHandleARFile(t *testing.T) {
	file, err := os.Open("testdata/test.deb")
	assert.Nil(t, err)
	defer file.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	handler := newARHandler()
	archiveChan, err := handler.HandleFile(context.AddLogger(ctx), file)
	assert.NoError(t, err)

	wantChunkCount := 102
	count := 0
	for range archiveChan {
		count++
	}

	assert.Equal(t, wantChunkCount, count)
}

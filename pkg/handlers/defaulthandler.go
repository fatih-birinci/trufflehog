package handlers

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/gabriel-vasile/mimetype"
	"github.com/mholt/archiver/v4"
	diskbufferreader "github.com/trufflesecurity/disk-buffer-reader"

	"github.com/trufflesecurity/trufflehog/v3/pkg/common"
	logContext "github.com/trufflesecurity/trufflehog/v3/pkg/context"
	"github.com/trufflesecurity/trufflehog/v3/pkg/sources"
)

type ctxKey int

const (
	depthKey          ctxKey = iota
	defaultBufferSize        = 512
)

var (
	maxDepth   = 5
	maxSize    = 250 * 1024 * 1024 // 250 MB
	maxTimeout = time.Duration(30) * time.Second
)

// SetArchiveMaxSize sets the maximum size of the archive.
func SetArchiveMaxSize(size int) { maxSize = size }

// SetArchiveMaxDepth sets the maximum depth of the archive.
func SetArchiveMaxDepth(depth int) { maxDepth = depth }

// SetArchiveMaxTimeout sets the maximum timeout for the archive handler.
func SetArchiveMaxTimeout(timeout time.Duration) { maxTimeout = timeout }

// defaultHandler provides a base implementation for file handlers, encapsulating common behaviors
// needed across different handlers. This handler is embedded in other specialized handlers to ensure
// consistent application of these common behaviors and to simplify the extension of handler functionalities.
type defaultHandler struct{}

// HandleFile processes the input as either an archive or non-archive based on its content,
// utilizing a single output channel. It first tries to identify the input as an archive. If it is an archive,
// it processes it accordingly; otherwise, it handles the input as non-archive content.
// The function returns a channel that will receive the extracted data bytes and an error if the initial setup fails.
func (h *defaultHandler) HandleFile(ctx logContext.Context, input *diskbufferreader.DiskBufferReader) (chan []byte, error) {
	// Shared channel for both archive and non-archive content.
	dataChan := make(chan []byte, defaultBufferSize)

	_, arReader, err := archiver.Identify("", input)
	if err != nil {
		if errors.Is(err, archiver.ErrNoMatch) {
			// Not an archive, handle as non-archive content in a separate goroutine.
			ctx.Logger().V(3).Info("File not recognized as an archive, handling as non-archive content.")
			go func() {
				defer close(dataChan)
				if err := h.handleNonArchiveContent(ctx, arReader, dataChan); err != nil {
					ctx.Logger().Error(err, "error handling non-archive content.")
				}
			}()
			return dataChan, nil
		}

		return nil, err
	}

	go func() {
		ctx, cancel := logContext.WithTimeout(ctx, maxTimeout)
		defer cancel()
		defer close(dataChan)

		if err := h.openArchive(ctx, 0, arReader, dataChan); err != nil {
			ctx.Logger().Error(err, "error unarchiving chunk.")
		}
	}()
	return dataChan, nil
}

var ErrMaxDepthReached = errors.New("max archive depth reached")

// openArchive recursively extracts content from an archive up to a maximum depth, handling nested archives if necessary.
// It takes a reader from which it attempts to identify and process the archive format. Depending on the archive type,
// it either decompresses or extracts the contents directly, sending data to the provided channel.
// Returns an error if the archive cannot be processed due to issues like exceeding maximum depth or unsupported formats.
func (h *defaultHandler) openArchive(ctx logContext.Context, depth int, reader io.Reader, archiveChan chan []byte) error {
	if common.IsDone(ctx) {
		return ctx.Err()
	}

	if depth > maxDepth {
		return ErrMaxDepthReached
	}

	format, arReader, err := archiver.Identify("", reader)
	switch {
	case err == nil:
		// Continue with the rest of the code.
	case errors.Is(err, archiver.ErrNoMatch):
		if depth > 0 {
			// If openArchive is called on an already extracted/decompressed file and the depth is greater than 0,
			// it means we are at least 1 layer deep in the archive. In this case, we should handle the content
			// as non-archive data by calling handleNonArchiveContent.
			return h.handleNonArchiveContent(ctx, arReader, archiveChan)
		}
		// If openArchive is called on the root (depth == 0) and we can't identify the format,
		// it means we can't handle the content at all. Return the archiver.ErrNoMatch error.
		return err
	default:
		// Some other error occurred.
		return fmt.Errorf("error identifying archive: %w", err)
	}

	switch archive := format.(type) {
	case archiver.Decompressor:
		// Decompress tha archive and feed the decompressed data back into the archive handler to extract any nested archives.
		compReader, err := archive.OpenReader(arReader)
		if err != nil {
			return err
		}

		defer compReader.Close()

		return h.openArchive(ctx, depth+1, compReader, archiveChan)
	case archiver.Extractor:
		return archive.Extract(logContext.WithValue(ctx, depthKey, depth+1), arReader, nil, h.extractorHandler(archiveChan))
	default:
		return fmt.Errorf("unknown archive type: %s", format.Name())
	}
}

// extractorHandler creates a closure that handles individual files extracted by an archiver.
// It logs the extraction, checks for cancellation, and decides whether to skip the file based on its name or type,
// particularly for binary files if configured to skip. If the file is not skipped, it recursively calls openArchive
// to handle nested archives or to continue processing based on the file's content and depth in the archive structure.
func (h *defaultHandler) extractorHandler(archiveChan chan []byte) func(context.Context, archiver.File) error {
	return func(ctx context.Context, file archiver.File) error {
		lCtx := logContext.WithValues(
			logContext.AddLogger(ctx),
			"filename", file.Name(),
			"size", file.Size(),
		)
		lCtx.Logger().V(5).Info("Handling extracted file.")

		if common.IsDone(ctx) {
			return ctx.Err()
		}

		depth := 0
		if ctxDepth, ok := ctx.Value(depthKey).(int); ok {
			depth = ctxDepth
		}
		if int(file.Size()) > maxSize {
			lCtx.Logger().V(3).Info("skipping file due to size")
			return nil
		}

		if common.SkipFile(file.Name()) {
			lCtx.Logger().V(5).Info("skipping file")
			return nil
		}

		fReader, err := file.Open()
		if err != nil {
			return err
		}
		defer fReader.Close()

		return h.openArchive(lCtx, depth, fReader, archiveChan)
	}
}

// handleNonArchiveContent processes files that do not contain nested archives, serving as the final stage in the
// extraction/decompression process. It reads the content to detect its MIME type and decides whether to skip based
// on the type, particularly for binary files. It manages reading file chunks and writing them to the archive channel,
// effectively collecting the final bytes for further processing. This function is a key component in ensuring that all
// file content, regardless of being an archive or not, is handled appropriately.
func (h *defaultHandler) handleNonArchiveContent(ctx logContext.Context, reader io.Reader, archiveChan chan []byte) error {
	bufReader := bufio.NewReaderSize(reader, defaultBufferSize)
	// A buffer of 512 bytes is used since many file formats store their magic numbers within the first 512 bytes.
	// If fewer bytes are read, MIME type detection may still succeed.
	buffer, err := bufReader.Peek(defaultBufferSize)
	if err != nil && !errors.Is(err, io.EOF) {
		return fmt.Errorf("unable to read file for MIME type detection: %w", err)
	}

	mime := mimetype.Detect(buffer)
	mimeT := mimeType(mime.String())

	if common.SkipFile(mime.Extension()) {
		ctx.Logger().V(5).Info("skipping file", "ext", mimeT)
		return nil
	}

	chunkReader := sources.NewChunkReader()
	chunkResChan := chunkReader(ctx, bufReader)
	for data := range chunkResChan {
		if err := data.Error(); err != nil {
			ctx.Logger().Error(err, "error reading chunk")
			continue
		}
		if err := common.CancellableWrite(ctx, archiveChan, data.Bytes()); err != nil {
			return err
		}
	}
	return nil
}

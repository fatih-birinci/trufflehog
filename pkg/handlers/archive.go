package handlers

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/mholt/archiver/v4"
	log "github.com/sirupsen/logrus"
)

var (
	maxDepth = 5
)

// Archive is a handler for extracting and decompressing archives.
type Archive struct {
	newChunkData bytes.Buffer
	maxSize      int
	size         int
}

// New sets a default maximum size and current size counter.
func (d *Archive) New() {
	d.maxSize = 20 * 1024 * 1024 // 20MB
	d.size = 0
}

// FromFile extracts the files from an archive.
func (d *Archive) FromFile(data io.Reader) chan ([]byte) {
	ctx := context.Background()
	archiveChan := make(chan ([]byte), 512)
	go func() {
		defer close(archiveChan)
		err := d.openArchive(ctx, 0, data, archiveChan)
		if err != nil {
			if errors.Is(err, archiver.ErrNoMatch) {
				return
			}
			log.WithError(err).Debug("Error unarchiving chunk.")
		}
	}()
	return archiveChan
}

// openArchive takes a reader and extracts the contents up to the maximum depth.
func (d *Archive) openArchive(ctx context.Context, depth int, reader io.Reader, archiveChan chan ([]byte)) error {
	if depth >= maxDepth {
		return fmt.Errorf("max archive depth reached")
	}
	format, reader, err := archiver.Identify("", reader)
	if err != nil {
		if errors.Is(err, archiver.ErrNoMatch) && depth > 0 {
			for {
				chunk := make([]byte, 10*1024)
				n, _ := reader.Read(chunk)
				archiveChan <- chunk
				if n < 512 {
					break
				}
			}
			fileBytes, err := d.ReadToMax(reader)
			if err != nil {
				return err
			}
			archiveChan <- fileBytes
			return nil
		}
		return err
	}
	switch archive := format.(type) {
	case archiver.Extractor:
		err := archive.Extract(context.WithValue(ctx, "depth", depth+1), reader, nil, d.extractorHandler(archiveChan))
		if err != nil {
			return err
		}
		return nil
	case archiver.Decompressor:
		compReader, err := archive.OpenReader(reader)
		if err != nil {
			return err
		}
		fileBytes, err := d.ReadToMax(compReader)
		if err != nil {
			return err
		}
		newReader := bytes.NewReader(fileBytes)
		return d.openArchive(ctx, depth+1, newReader, archiveChan)
	}
	return fmt.Errorf("Unknown archive type: %s", format.Name())
}

// IsFiletype returns true if the provided reader is an archive.
func (d *Archive) IsFiletype(reader io.Reader) bool {
	format, readerB, err := archiver.Identify("", reader)
	reader = readerB
	if err != nil {
		return false
	}
	switch format.(type) {
	case archiver.Extractor:
		return true
	case archiver.Decompressor:
		return true
	}
	return false
}

// extractorHandler is applied to each file in an archiver.Extractor file.
func (d *Archive) extractorHandler(archiveChan chan ([]byte)) func(context.Context, archiver.File) error {
	return func(ctx context.Context, f archiver.File) error {
		log.WithField("filename", f.Name()).Trace("Handling extracted file.")
		depth := 0
		if ctxDepth, ok := ctx.Value("depth").(int); ok {
			depth = ctxDepth
		}

		fReader, err := f.Open()
		if err != nil {
			return err
		}
		fileBytes, err := d.ReadToMax(fReader)
		if err != nil {
			return err
		}
		fileContent := bytes.NewReader(fileBytes)

		err = d.openArchive(ctx, depth, fileContent, archiveChan)
		if err != nil {
			return err
		}
		return nil
	}
}

// ReadToMax reads up to the max size.
func (d *Archive) ReadToMax(reader io.Reader) ([]byte, error) {
	fileContent := bytes.Buffer{}
	log.Tracef("Remaining buffer capacity: %d", d.maxSize-d.size)
	for i := 0; i <= d.maxSize/512; i++ {
		fileChunk := make([]byte, 512)
		bRead, err := reader.Read(fileChunk)
		if err != nil && !errors.Is(err, io.EOF) {
			return []byte{}, err
		}
		d.size += bRead
		fileContent.Write(fileChunk[0:bRead])
		if bRead < 512 {
			break
		}
		if d.size >= d.maxSize && bRead == 512 {
			log.Debug("Max archive size reached.")
			break
		}
	}
	return fileContent.Bytes(), nil
}

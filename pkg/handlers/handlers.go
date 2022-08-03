package handlers

import (
	"io"

	"github.com/trufflesecurity/trufflehog/v3/pkg/sources"
)

func DefaultHandlers() []Handler {
	return []Handler{
		&Archive{},
	}
}

type Handler interface {
	FromFile(io.Reader) chan ([]byte)
	IsFiletype(io.Reader) (io.Reader, bool)
	New()
}

func HandleFile(file io.Reader, chunkSkel *sources.Chunk, chunksChan chan (*sources.Chunk)) bool {
	for _, handler := range DefaultHandlers() {
		handler.New()
		var isType bool
		file, isType = handler.IsFiletype(file)
		if !isType {
			continue
		}
		handlerChan := handler.FromFile(file)
		for data := range handlerChan {
			chunk := *chunkSkel
			chunk.Data = data
			chunksChan <- &chunk
		}
		return true
	}
	return false
}

package sources

import (
	"fmt"
	"testing"

	"github.com/trufflesecurity/trufflehog/v3/pkg/context"
	"github.com/trufflesecurity/trufflehog/v3/pkg/pb/sourcespb"
	"google.golang.org/protobuf/types/known/anypb"
)

// DummySource implements Source and is used for testing a SourceManager.
type DummySource struct {
	sourceID int64
	jobID    int64
	chunker
}

func (d *DummySource) Type() sourcespb.SourceType { return 1337 }
func (d *DummySource) SourceID() int64            { return d.sourceID }
func (d *DummySource) JobID() int64               { return d.jobID }
func (d *DummySource) Init(_ context.Context, _ string, jobID, sourceID int64, _ bool, _ *anypb.Any, _ int) error {
	d.sourceID = sourceID
	d.jobID = jobID
	return nil
}
func (d *DummySource) GetProgress() *Progress { return nil }

// Interface to easily test different chunking methods.
type chunker interface {
	Chunks(context.Context, chan *Chunk) error
	ChunkUnit(ctx context.Context, unit SourceUnit, reporter ChunkReporter) error
	Enumerate(ctx context.Context, reporter UnitReporter) error
}

// Chunk method that writes count bytes to the channel before returning.
type counterChunker struct {
	chunkCounter byte
	count        int
}

func (c *counterChunker) Chunks(_ context.Context, ch chan *Chunk) error {
	for i := 0; i < c.count; i++ {
		ch <- &Chunk{Data: []byte{c.chunkCounter}}
		c.chunkCounter++
	}
	return nil
}

// countChunk implements SourceUnit.
type countChunk byte

func (c countChunk) SourceUnitID() string { return fmt.Sprintf("countChunk(%d)", c) }

func (c *counterChunker) Enumerate(ctx context.Context, reporter UnitReporter) error {
	for i := 0; i < c.count; i++ {
		if err := reporter.UnitOk(ctx, countChunk(byte(i))); err != nil {
			return err
		}
	}
	return nil
}

func (c *counterChunker) ChunkUnit(ctx context.Context, unit SourceUnit, reporter ChunkReporter) error {
	return reporter.ChunkOk(ctx, Chunk{Data: []byte{byte(unit.(countChunk))}})
}

// Chunk method that always returns an error.
type errorChunker struct{ error }

func (c errorChunker) Chunks(context.Context, chan *Chunk) error                  { return c }
func (c errorChunker) Enumerate(context.Context, UnitReporter) error              { return c }
func (c errorChunker) ChunkUnit(context.Context, SourceUnit, ChunkReporter) error { return c }

// enrollDummy is a helper function to enroll a DummySource with a SourceManager.
func enrollDummy(mgr *SourceManager, chunkMethod chunker) (handle, error) {
	return mgr.Enroll(context.Background(), "dummy", 1337,
		func(ctx context.Context, jobID, sourceID int64) (Source, error) {
			source := &DummySource{chunker: chunkMethod}
			if err := source.Init(ctx, "dummy", jobID, sourceID, true, nil, 42); err != nil {
				return nil, err
			}
			return source, nil
		})
}

// tryRead is a helper function that will try to read from a channel and return
// an error if it cannot.
func tryRead(ch <-chan *Chunk) (*Chunk, error) {
	select {
	case chunk := <-ch:
		return chunk, nil
	default:
		return nil, fmt.Errorf("no chunk available")
	}
}

func TestSourceManagerRun(t *testing.T) {
	mgr := NewManager(WithBufferedOutput(8))
	handle, err := enrollDummy(mgr, &counterChunker{count: 1})
	if err != nil {
		t.Fatalf("unexpected error enrolling source: %v", err)
	}
	for i := 0; i < 3; i++ {
		if err := mgr.Run(context.Background(), handle); err != nil {
			t.Fatalf("unexpected error running source: %v", err)
		}
		chunk, err := tryRead(mgr.Chunks())
		if err != nil {
			t.Fatalf("reading chunk failed: %v", err)
		}
		if chunk.Data[0] != byte(i) {
			t.Fatalf("unexpected chunk value, wanted %v, got: %v", chunk.Data[0], i)
		}

		// The Chunks channel should be empty now.
		if chunk, err := tryRead(mgr.Chunks()); err == nil {
			t.Fatalf("unexpected chunk found: %+v", chunk)
		}
	}
}

func TestSourceManagerWait(t *testing.T) {
	mgr := NewManager()
	handle, err := enrollDummy(mgr, &counterChunker{count: 1})
	if err != nil {
		t.Fatalf("unexpected error enrolling source: %v", err)
	}
	// Asynchronously run the source.
	if err := mgr.ScheduleRun(context.Background(), handle); err != nil {
		t.Fatalf("unexpected error scheduling run: %v", err)
	}
	// Read the 1 chunk we're expecting so Waiting completes.
	<-mgr.Chunks()
	// Wait for all resources to complete.
	if err := mgr.Wait(); err != nil {
		t.Fatalf("unexpected error waiting: %v", err)
	}
	// Enroll and run should return an error now.
	if _, err := enrollDummy(mgr, &counterChunker{count: 1}); err == nil {
		t.Fatalf("expected enroll to fail")
	}
	if err := mgr.ScheduleRun(context.Background(), handle); err == nil {
		t.Fatalf("expected scheduling run to fail")
	}
}

func TestSourceManagerError(t *testing.T) {
	mgr := NewManager()
	handle, err := enrollDummy(mgr, errorChunker{fmt.Errorf("oops")})
	if err != nil {
		t.Fatalf("unexpected error enrolling source: %v", err)
	}
	// A synchronous run should fail.
	if err := mgr.Run(context.Background(), handle); err == nil {
		t.Fatalf("expected run to fail")
	}
	// Scheduling a run should not fail, but the error should surface in
	// Wait().
	if err := mgr.ScheduleRun(context.Background(), handle); err != nil {
		t.Fatalf("unexpected error scheduling run: %v", err)
	}
	if err := mgr.Wait(); err == nil {
		t.Fatalf("expected wait to fail")
	}
}

func TestSourceManagerReport(t *testing.T) {
	for _, opts := range [][]func(*SourceManager){
		{WithBufferedOutput(8)},
		{WithBufferedOutput(8), WithSourceUnits()},
		{WithBufferedOutput(8), WithSourceUnits(), WithConcurrentUnits(1)},
	} {
		mgr := NewManager(opts...)
		handle, err := enrollDummy(mgr, &counterChunker{count: 4})
		if err != nil {
			t.Fatalf("unexpected error enrolling source: %v", err)
		}
		// Synchronously run the source.
		if err := mgr.Run(context.Background(), handle); err != nil {
			t.Fatalf("unexpected error running source: %v", err)
		}
		report := mgr.Report(handle)
		if report == nil {
			t.Fatalf("expected a report")
		}
		if err := report.Errors(); err != nil {
			t.Fatalf("unexpected error in report: %v", err)
		}
		if report.TotalChunks != 4 {
			t.Fatalf("expected report to have 4 chunks, got: %d", report.TotalChunks)
		}
	}
}

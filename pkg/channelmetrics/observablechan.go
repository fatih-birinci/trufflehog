// Package channelmetrics provides a flexible way to wrap Go channels with
// additional metrics collection capabilities. This allows for monitoring
// and tracking of channel usage and performance using different metrics backends.
package channelmetrics

import (
	"time"

	"github.com/trufflesecurity/trufflehog/v3/pkg/common"
	"github.com/trufflesecurity/trufflehog/v3/pkg/context"
)

// MetricsCollector is an interface for collecting metrics. Implementations
// of this interface can be used to record various channel metrics.
type MetricsCollector interface {
	RecordProduceDuration(duration time.Duration)
	RecordConsumeDuration(duration time.Duration)
	RecordChannelLen(size int)
	RecordChannelCap(capacity int)
}

// ObservableChan wraps a Go channel and collects metrics about its usage.
// It supports any type of channel and records metrics using a provided
// MetricsCollector implementation.
type ObservableChan[T any] struct {
	ch        chan T
	metrics   MetricsCollector
	bufferCap int
}

// NewObservableChan creates a new ObservableChan wrapping the provided channel.
// It records the channel's capacity immediately and sets up metrics collection
// using the provided MetricsCollector and channel name. The chanName is used to
// distinguish between metrics for different channels by incorporating it into
// the metric names.
func NewObservableChan[T any](ch chan T, metrics MetricsCollector) *ObservableChan[T] {
	oChan := &ObservableChan[T]{
		ch:        ch,
		metrics:   metrics,
		bufferCap: cap(ch),
	}
	oChan.RecordChannelCapacity() // Record capacity immediately
	return oChan
}

// Close closes the channel and records the current size of the channel buffer.
func (oc *ObservableChan[T]) Close() {
	close(oc.ch)
	oc.RecordChannelLen()
}

// Send sends an item into the channel and records the duration taken to do so.
// It also updates the current size of the channel buffer.
func (oc *ObservableChan[T]) Send(ctx context.Context, item T) {
	startTime := time.Now()
	defer func() {
		oc.metrics.RecordProduceDuration(time.Since(startTime))
		oc.RecordChannelLen()
	}()
	if err := common.CancellableWrite(ctx, oc.ch, item); err != nil {
		ctx.Logger().Error(err, "failed to write item to observable channel")
	}
}

// Recv receives an item from the channel and records the duration taken to do so.
// It also updates the current size of the channel buffer.
func (oc *ObservableChan[T]) Recv(_ context.Context) T {
	startTime := time.Now()
	defer func() {
		oc.metrics.RecordConsumeDuration(time.Since(startTime))
		oc.RecordChannelLen()
	}()
	return <-oc.ch
}

// RecordChannelCapacity records the capacity of the channel buffer.
func (oc *ObservableChan[T]) RecordChannelCapacity() {
	oc.metrics.RecordChannelCap(oc.bufferCap)
}

// RecordChannelLen records the current size of the channel buffer.
func (oc *ObservableChan[T]) RecordChannelLen() {
	oc.metrics.RecordChannelLen(len(oc.ch))
}

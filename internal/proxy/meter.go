package proxy

import (
	"io"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

type AppMetrics struct {
	BytesIn         uint64 `json:"bytes_in"`
	BytesOut        uint64 `json:"bytes_out"`
	TotalRequests   uint64 `json:"total_requests"`
	ActiveRequests  int64  `json:"active_requests"`
	BytesPerSecIn   uint64 `json:"bytes_per_sec_in"`
	BytesPerSecOut  uint64 `json:"bytes_per_sec_out"`
	LastSampledTime time.Time
	lastSampleIn    uint64
	lastSampleOut   uint64
}

type BandwidthMeter struct {
	mu         sync.RWMutex
	apps       map[string]*AppMetrics
	totalIn    uint64
	totalOut   uint64
	totalReqs  uint64
	activeReqs int64
}

func NewBandwidthMeter() *BandwidthMeter {
	bm := &BandwidthMeter{
		apps: make(map[string]*AppMetrics),
	}
	go bm.rateSamplerLoop()
	return bm
}

func (m *BandwidthMeter) getOrCreateApp(appID string) *AppMetrics {
	m.mu.Lock()
	defer m.mu.Unlock()

	metrics, exists := m.apps[appID]
	if !exists {
		metrics = &AppMetrics{LastSampledTime: time.Now()}
		m.apps[appID] = metrics
	}
	return metrics
}

func (m *BandwidthMeter) RecordRequestStart(appID string) {
	atomic.AddInt64(&m.activeReqs, 1)
	atomic.AddUint64(&m.totalReqs, 1)
	if appID != "" {
		app := m.getOrCreateApp(appID)
		atomic.AddInt64(&app.ActiveRequests, 1)
		atomic.AddUint64(&app.TotalRequests, 1)
	}
}

func (m *BandwidthMeter) RecordRequestEnd(appID string) {
	atomic.AddInt64(&m.activeReqs, -1)
	if appID != "" {
		app := m.getOrCreateApp(appID)
		atomic.AddInt64(&app.ActiveRequests, -1)
	}
}

func (m *BandwidthMeter) RecordBytesIn(appID string, n int64) {
	if n <= 0 {
		return
	}
	u := uint64(n)
	atomic.AddUint64(&m.totalIn, u)
	if appID != "" {
		app := m.getOrCreateApp(appID)
		atomic.AddUint64(&app.BytesIn, u)
	}
}

func (m *BandwidthMeter) RecordBytesOut(appID string, n int64) {
	if n <= 0 {
		return
	}
	u := uint64(n)
	atomic.AddUint64(&m.totalOut, u)
	if appID != "" {
		app := m.getOrCreateApp(appID)
		atomic.AddUint64(&app.BytesOut, u)
	}
}

func (m *BandwidthMeter) GetAppStats(appID string) AppMetrics {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if app, ok := m.apps[appID]; ok {
		return AppMetrics{
			BytesIn:         atomic.LoadUint64(&app.BytesIn),
			BytesOut:        atomic.LoadUint64(&app.BytesOut),
			TotalRequests:   atomic.LoadUint64(&app.TotalRequests),
			ActiveRequests:  atomic.LoadInt64(&app.ActiveRequests),
			BytesPerSecIn:   atomic.LoadUint64(&app.BytesPerSecIn),
			BytesPerSecOut:  atomic.LoadUint64(&app.BytesPerSecOut),
			LastSampledTime: app.LastSampledTime,
		}
	}
	return AppMetrics{}
}

func (m *BandwidthMeter) rateSamplerLoop() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		m.mu.Lock()
		now := time.Now()
		for _, app := range m.apps {
			curIn := atomic.LoadUint64(&app.BytesIn)
			curOut := atomic.LoadUint64(&app.BytesOut)

			deltaIn := curIn - app.lastSampleIn
			deltaOut := curOut - app.lastSampleOut

			atomic.StoreUint64(&app.BytesPerSecIn, deltaIn)
			atomic.StoreUint64(&app.BytesPerSecOut, deltaOut)

			app.lastSampleIn = curIn
			app.lastSampleOut = curOut
			app.LastSampledTime = now
		}
		m.mu.Unlock()
	}
}

// CountingResponseWriter wraps http.ResponseWriter to track egress bandwidth.
type CountingResponseWriter struct {
	http.ResponseWriter
	appID   string
	meter   *BandwidthMeter
	written int64
	code    int
}

func NewCountingResponseWriter(w http.ResponseWriter, appID string, meter *BandwidthMeter) *CountingResponseWriter {
	return &CountingResponseWriter{
		ResponseWriter: w,
		appID:          appID,
		meter:          meter,
		code:           http.StatusOK,
	}
}

func (rw *CountingResponseWriter) WriteHeader(statusCode int) {
	rw.code = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

func (rw *CountingResponseWriter) Write(b []byte) (int, error) {
	n, err := rw.ResponseWriter.Write(b)
	rw.written += int64(n)
	rw.meter.RecordBytesOut(rw.appID, int64(n))
	return n, err
}

func (rw *CountingResponseWriter) Flush() {
	if f, ok := rw.ResponseWriter.(http.Flusher); ok {
		f.Flush()
	}
}

// CountingReader wraps io.Reader to track ingress bandwidth.
type CountingReader struct {
	io.ReadCloser
	appID string
	meter *BandwidthMeter
}

func NewCountingReader(r io.ReadCloser, appID string, meter *BandwidthMeter) *CountingReader {
	return &CountingReader{ReadCloser: r, appID: appID, meter: meter}
}

func (cr *CountingReader) Read(p []byte) (int, error) {
	n, err := cr.ReadCloser.Read(p)
	if n > 0 {
		cr.meter.RecordBytesIn(cr.appID, int64(n))
	}
	return n, err
}

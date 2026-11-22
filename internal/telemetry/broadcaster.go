package telemetry

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type TelemetryBroadcaster struct {
	collector *TelemetryCollector
}

func NewTelemetryBroadcaster(collector *TelemetryCollector) *TelemetryBroadcaster {
	return &TelemetryBroadcaster{collector: collector}
}

func (tb *TelemetryBroadcaster) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/telemetry/metrics", tb.handleGetMetrics)
	mux.HandleFunc("GET /api/v1/telemetry/stream", tb.handleStreamMetrics)
}

func (tb *TelemetryBroadcaster) handleGetMetrics(w http.ResponseWriter, r *http.Request) {
	metrics := tb.collector.Collect()
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(metrics)
}

func (tb *TelemetryBroadcaster) handleStreamMetrics(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming unsupported", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Accel-Buffering", "no")

	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	// Initial emission
	initial := tb.collector.Collect()
	initialJSON, _ := json.Marshal(initial)
	_, _ = fmt.Fprintf(w, "data: %s\n\n", initialJSON)
	flusher.Flush()

	for {
		select {
		case <-r.Context().Done():
			return
		case <-ticker.C:
			m := tb.collector.Collect()
			data, err := json.Marshal(m)
			if err != nil {
				continue
			}
			_, _ = fmt.Fprintf(w, "data: %s\n\n", data)
			flusher.Flush()
		}
	}
}

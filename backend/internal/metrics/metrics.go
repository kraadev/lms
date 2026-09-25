package metrics

import (
	"fmt"
	"net/http"
	"sync/atomic"
)

var (
	TotalRequests   uint64
	ActiveWSCount   int64
)

func Handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; version=0.0.4")
	fmt.Fprintf(w, "# HELP lms_http_requests_total Total number of HTTP requests processed.\n")
	fmt.Fprintf(w, "# TYPE lms_http_requests_total counter\n")
	fmt.Fprintf(w, "lms_http_requests_total %d\n\n", atomic.LoadUint64(&TotalRequests))
	fmt.Fprintf(w, "# HELP lms_active_websockets Current number of active WebSocket connections.\n")
	fmt.Fprintf(w, "# TYPE lms_active_websockets gauge\n")
	fmt.Fprintf(w, "lms_active_websockets %d\n", atomic.LoadInt64(&ActiveWSCount))
}

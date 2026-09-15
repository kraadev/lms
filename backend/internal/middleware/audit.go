package middleware

import (
	"log"
	"net/http"
	"time"
)

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}

// AuditLogger logs state-modifying actions (POST, PUT, PATCH, DELETE)
func AuditLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}

		next.ServeHTTP(rec, r)

		if r.Method == "POST" || r.Method == "PUT" || r.Method == "PATCH" || r.Method == "DELETE" {
			duration := time.Since(start).Milliseconds()
			log.Printf("[AUDIT] method=%s path=%s status=%d duration_ms=%d ip=%s",
				r.Method, r.URL.Path, rec.status, duration, r.RemoteAddr)
		}
	})
}

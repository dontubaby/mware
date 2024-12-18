package middleware

import (
	"context"
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"time"
)

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (rec *statusRecorder) WriteHeader(code int) {
	rec.status = code
	rec.ResponseWriter.WriteHeader(code)
}

const digitCount = 12

func ReuestIDGenerator() string {
	var result []byte
	for i := 0; i < digitCount; i++ {
		num, _ := rand.Int(rand.Reader, big.NewInt(10))
		result = append(result, byte('0'+num.Int64()))
	}
	return string(result)
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := statusRecorder{ResponseWriter: w}
		next.ServeHTTP(&rec, r)
		log.Printf("IPv6: %s, request_id_new: %s, status: %d, duration: %s",
			r.RemoteAddr,
			r.Context().Value("request_id"),
			rec.status,
			time.Since(start))
	})
}

func RequstIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var ctx context.Context
		requestID := r.URL.Query().Get("request_id")
		if requestID == "" {
			requestID = ReuestIDGenerator()
		}
		fmt.Println("Request id BEFORE", r.Context().Value("request_id"))
		if r.Context().Value("request_id") == nil {
			ctx = context.WithValue(r.Context(), "request_id", requestID)
		}
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

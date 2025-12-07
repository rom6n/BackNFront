package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/rom6n/random-greetings"
)

var (
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "otello_http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"path", "method", "code"},
	)

	httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "otello_http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"path", "method"},
	)

	fileWritesTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "otello_file_writes_total",
			Help: "Total number of successful file write attempts",
		},
	)

	fileWriteErrorsTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "otello_file_write_errors_total",
			Help: "Total number of file write errors",
		},
	)
)

func init() {
	// Регистрируем метрики
	prometheus.MustRegister(httpRequestsTotal, httpRequestDuration, fileWritesTotal, fileWriteErrorsTotal)
}

func instrumentHandler(path string, h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// простой wrapper чтобы получить код — используем ResponseWriter wrapper
		ww := &statusRecordingResponseWriter{ResponseWriter: w, status: 200}
		h(ww, r)

		duration := time.Since(start).Seconds()
		httpRequestDuration.WithLabelValues(path, r.Method).Observe(duration)
		httpRequestsTotal.WithLabelValues(path, r.Method, http.StatusText(ww.status)).Inc()
	}
}

type statusRecordingResponseWriter struct {
	http.ResponseWriter
	status int
}

func (w *statusRecordingResponseWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

func main() {
	// handler, который пишет в файл index.txt
	http.HandleFunc("/", instrumentHandler("/", func(w http.ResponseWriter, _ *http.Request) {
		f, err := os.OpenFile("index.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			fileWriteErrorsTotal.Inc()
			log.Println("open file error:", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
		defer f.Close()

		if _, err := f.WriteString("Hello World 🐋🐋🐋\n"); err != nil {
			fileWriteErrorsTotal.Inc()
			log.Println("write file error:", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
		fileWritesTotal.Inc()

		w.Write([]byte("Hello World, file opened 🐋"))
	}))

	// prometheus handler
	http.Handle("/metrics", promhttp.Handler())

	http.Handle("/health", instrumentHandler("/health", func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte(fmt.Sprintf("OK - %v %v\n", greetings.GetRandomGreeting(), req.Header.Get("X-Real-IP"))))
	}))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Println("listening on :", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

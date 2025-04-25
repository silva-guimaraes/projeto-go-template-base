package logging

import (
	"errors"
	"fmt"
	"foobar/metrics"
	"log"
	"net/http"
	"runtime/debug"
	"time"
)

type wrappedWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *wrappedWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
	w.statusCode = statusCode
}

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		wrapped := &wrappedWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		var panicked bool
		err := recoverServeHTTP(next, wrapped, r)
		if err != nil {
			wrapped.Write([]byte("500 error"))
			log.Println(err)
			panicked = true
		}

		var (
			latency    = time.Since(start)
			path       = r.URL.Path
			statusCode = wrapped.statusCode
		)
		if path == "/metrics" && statusCode == 200 {
			return
		}
		log.Printf("%v %v %v %v", wrapped.statusCode, r.Method, r.URL.Path, latency)

		metrics.ResponseLatencyObserve(latency, statusCode, panicked, r)
	})
}

func recoverServeHTTP(next http.Handler, w http.ResponseWriter, r *http.Request) (e error) {
	defer func() {
		if rec := recover(); rec != nil {
			var err error
			switch v := rec.(type) {
			case error:
				err = v
			default:
				err = fmt.Errorf("%v", v)
			}
			stack := debug.Stack()
			w.WriteHeader(http.StatusInternalServerError)
			e = errors.Join(err, fmt.Errorf(string(stack)))
		}
	}()
	next.ServeHTTP(w, r)
	return
}

func Error(err error) {
	if err == nil {
		log.Printf("LOG: Erro Nulo!?\n")
	} else {
		log.Printf("LOG: %s\n", err.Error())
	}
	log.Println(string(debug.Stack()))
}

func InternalError(w http.ResponseWriter, err error) {
	Error(err)
	http.Error(w, `500`, http.StatusInternalServerError)
}

func FormValueMissing(field string, r *http.Request) error {
	return fmt.Errorf(
		"campo '%s' em formulário não pôde ser encontrado. Content-Type: %s",
		field,
		r.Header.Get("Content-Type"),
	)
}

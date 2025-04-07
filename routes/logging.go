package routes

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type wrappedWriter struct {
	http.ResponseWriter
	statusCode int
}

var (
	opsProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "myapp_processed_ops_total",
		Help: "The total number of processed events",
	})
)

func (w *wrappedWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
	w.statusCode = statusCode
}

func logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		wrapped := &wrappedWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		err := recoverServeHTTP(next, wrapped, r)
		if err != nil {
			log.Println(err)
		}

		log.Printf("%v %v %v %v", wrapped.statusCode, r.Method, r.URL.Path, time.Since(start))
		opsProcessed.Inc()
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

func logError(err error) {
	if err == nil {
		log.Printf("LOG: Erro Nulo!?\n")
	} else {
		log.Printf("LOG: %s\n", err.Error())
	}
	log.Println(string(debug.Stack()))
}

func logInternalError(w http.ResponseWriter, err error) {
	logError(err)
	http.Error(w, `500`, http.StatusInternalServerError)
}

func httpErrorf(w http.ResponseWriter, status int, formatString string, format ...interface{}) {
	http.Error(w, fmt.Sprintf(formatString, format...), status)
}

func formValueMissing(field string, r *http.Request) error {
	return fmt.Errorf(
		"campo '%s' em formulário não pôde ser encontrado. Content-Type: %s",
		field,
		r.Header.Get("Content-Type"),
	)
}

package routes

import (
	"fmt"
	"log"
	"net/http"
	"os"
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

func logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		wrapped := &wrappedWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		next.ServeHTTP(wrapped, r)

		log.Printf("%v %v %v %v", wrapped.statusCode, r.Method, r.URL.Path, time.Since(start))
	})
}

func logPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				trace := fmt.Sprintf("Catastrophique!\n%v\n%s", err, string(debug.Stack()))
				log.Println(trace)
				w.WriteHeader(500)
				if os.Getenv("PROD") == "" {
					w.Write([]byte(trace))
				}
			}
		}()
		next.ServeHTTP(w, r)
	})
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

func formValueMissing(field string, w http.ResponseWriter, r *http.Request) error {
	return fmt.Errorf(
		"campo '%s' em formulário não pôde ser encontrado. Content-Type: %s",
		field,
		r.Header.Get("Content-Type"),
	)
}

// func logFormValueMissing(field string, w http.ResponseWriter, r *http.Request) {
// 	logInternalError(w, formatFormValueMissing(field, w, r))
// }

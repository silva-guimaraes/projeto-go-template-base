package index

import (
	"foobar/model"
	"net/http"

	"github.com/a-h/templ"
)

func Index(w http.ResponseWriter, r *http.Request, _ *model.User) (templ.Component, error) {
	if r.URL.Path != "/index.html" && r.URL.Path != "/" || r.Method != "GET" {
		return notFound(w)
	}
	return index(), nil
}

func notFound(w http.ResponseWriter) (templ.Component, error) {
	w.WriteHeader(404)
	return fourOfour(), nil
}

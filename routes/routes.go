package routes

import (
	"foobar/routes/auth"
	"foobar/routes/logging"
	. "foobar/routes/logging"
	. "foobar/routes/pages/index"
	. "foobar/routes/pages/logged"
	. "foobar/routes/pages/login"
	. "foobar/routes/pages/sobre"
	"net/http"
	"os"
	"path/filepath"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	staticDir string = filepath.Join("views", "static")
	mux       *http.ServeMux
)

func init() {
	mux = http.NewServeMux()
	if err := os.MkdirAll(staticDir, 777); err != nil {
		panic(err)
	}
	mux.Handle(
		"/static/",
		http.StripPrefix("/static/", http.FileServerFS(os.DirFS(staticDir))),
	)
	mux.HandleFunc(
		"/",
		getRouteMiddleware(Index),
	)
	mux.HandleFunc(
		"GET /login",
		getRouteMiddleware(Login),
	)
	mux.HandleFunc(
		"POST /login",
		redirectHtmxFormMiddleware(LogInPOST),
	)
	mux.HandleFunc(
		"GET /logged",
		getRouteMiddleware(Logged),
	)
	mux.HandleFunc(
		"GET /logout",
		logout,
	)
	mux.HandleFunc(
		"GET /sobre",
		getRouteMiddleware(Sobre),
	)
	mux.Handle(
		"/metrics",
		promhttp.Handler(),
	)
	mux.HandleFunc("/panic", func(w http.ResponseWriter, r *http.Request) {
		panic("não implementado")
	})
}

func Mux() http.Handler {
	return Logging(mux)
}

func logout(w http.ResponseWriter, r *http.Request) {
	session, err := auth.Store.Get(r, auth.SessionIdCookie)
	if err != nil {
		logging.InternalError(w, err)
		return
	}
	session.Options.MaxAge = -1
	if err = session.Save(r, w); err != nil {
		logging.InternalError(w, err)
		return
	}
	http.Redirect(w, r, "/login", http.StatusFound)
}

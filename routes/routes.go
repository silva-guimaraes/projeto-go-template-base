package routes

import (
	"context"
	"database/sql"
	"errors"
	"foobar/database"
	"foobar/views"
	"net/http"
	"os"
	"path/filepath"

	"golang.org/x/crypto/bcrypt"
)

var (
	staticDir = filepath.Join("views", "static")
)

func Register() http.Handler {
	mux := http.NewServeMux()

	if err := os.MkdirAll(staticDir, 777); err != nil {
		panic(err)
	}
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServerFS(os.DirFS(staticDir))))

	mux.HandleFunc("/", index)
	mux.HandleFunc("GET /login", logIn)
	mux.HandleFunc("POST /login", redirectHtmxFormMiddleware(logInPOST))
	mux.HandleFunc("GET /logged", logged)
	mux.HandleFunc("GET /logout", logout)
	mux.HandleFunc("GET /sobre", sobre)
	return logPanic(logging(mux))
}

func notFound(w http.ResponseWriter) {
	w.WriteHeader(404)
	if err := views.FourOfour().Render(context.Background(), w); err != nil {
		logInternalError(w, err)
		return
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/index.html" && r.URL.Path != "/" || r.Method != "GET" {
		notFound(w)
		return
	}
	if err := views.Index().Render(context.Background(), w); err != nil {
		logInternalError(w, err)
		return
	}
}

func sobre(w http.ResponseWriter, r *http.Request) {
	if err := views.Sobre().Render(context.Background(), w); err != nil {
		logInternalError(w, err)
		return
	}
}

func logIn(w http.ResponseWriter, r *http.Request) {

	_, err := currentUser(r)

	userLogged := err == nil

	if userLogged {
		http.Redirect(w, r, "/logged", http.StatusFound)
		return
	}

	if err := views.Login().Render(context.Background(), w); err != nil {
		logInternalError(w, err)
		return
	}
}

func logInPOST(w http.ResponseWriter, r *http.Request) (redirectURL, error) {

	email := r.FormValue("email")
	if email == "" {
		return "", formValueMissing("email", r)
	}

	senha := r.FormValue("senha")
	if senha == "" {
		return "", formValueMissing("senha", r)
	}

	tx := database.MustBeginTx()
	defer tx.Rollback()

	usuario, err := database.VerifyUser(tx, email, senha)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) || errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return "", newUserError(err, "senha ou email não conferem")
		} else {
			return "", err
		}
	}

	session, err := store.Get(r, sessionIdCookie)
	if err != nil {
		return "", err
	}

	userId, err := usuario.Id()
	if err != nil {
		return "", err
	}

	session.Values["userId"] = userId
	session.Options.SameSite = http.SameSiteDefaultMode

	if err := session.Save(r, w); err != nil {
		return "", err
	}

	return "/logged", nil
}

func logged(w http.ResponseWriter, r *http.Request) {
	usuario, err := currentUser(r)
	if err != nil {
		logError(err)
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Não Autorizado."))
		return
	}
	if err = views.Logged(usuario).Render(context.Background(), w); err != nil {
		logInternalError(w, err)
		return
	}
}

func logout(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, sessionIdCookie)
	if err != nil {
		logInternalError(w, err)
		return
	}
	session.Options.MaxAge = -1
	if err = session.Save(r, w); err != nil {
		logInternalError(w, err)
		return
	}
	http.Redirect(w, r, "/login", http.StatusFound)
}

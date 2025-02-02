package routes

import (
	"database/sql"
	"errors"
	"foobar/database"
	"html/template"
	"net/http"
	"os"
	"path/filepath"

	"golang.org/x/crypto/bcrypt"
)

var (
	viewsDir    = filepath.Join("views", "html")
	indexLayout = filepath.Join(viewsDir, "indexLayout.html")
	staticDir   = filepath.Join("views", "static")
)

func Register() http.Handler {
	mux := http.NewServeMux()

	if err := os.MkdirAll(staticDir, 777); err != nil {
		panic(err)
	}
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServerFS(os.DirFS(staticDir))))

	mux.HandleFunc("/", index)
	mux.HandleFunc("GET /login", logIn)
	mux.HandleFunc("POST /login", logInPOST)
	mux.HandleFunc("GET /logged", logged)
	mux.HandleFunc("GET /logout", logout)
	return logPanic(logging(mux))
}

func notFound(w http.ResponseWriter) {
	t, err := template.ParseFiles(indexLayout, filepath.Join(viewsDir, "404.html"))
	if err != nil {
		logInternalError(w, err)
		return
	}
	w.WriteHeader(404)
	if err := t.Execute(w, nil); err != nil {
		logInternalError(w, err)
		return
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/index.html" && r.URL.Path != "/" || r.Method != "GET" {
		notFound(w)
		return
	}
	t, err := template.ParseFiles(indexLayout, filepath.Join(viewsDir, "index.html"))
	if err != nil {
		logInternalError(w, err)
		return
	}
	if err := t.Execute(w, nil); err != nil {
		logInternalError(w, err)
		return
	}
}

func logIn(w http.ResponseWriter, r *http.Request) {

	_, err := currentUser(r)

	userLooged := err == nil

	if userLooged {
		http.Redirect(w, r, "/logged", http.StatusFound)
		return
	} else {
		t, err := template.ParseFiles(indexLayout, filepath.Join(viewsDir, "login.html"))
		if err != nil {
			logInternalError(w, err)
			return
		}
		if err := t.Execute(w, nil); err != nil {
			logInternalError(w, err)
			return
		}
	}
}

func logInPOST(w http.ResponseWriter, r *http.Request) {

	email := r.FormValue("email")
	if email == "" {
		formValueMissing("email", w, r)
		return
	}

	senha := r.FormValue("senha")
	if senha == "" {
		formValueMissing("senha", w, r)
		return
	}

	tx := database.MustBeginTx()
	defer tx.Rollback()

	usuario, err := database.VerifyUser(tx, email, senha)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) || errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			w.WriteHeader(422)
			w.Write([]byte(`senha ou email não conferem`))
			return
		} else {
			logInternalError(w, err)
			return
		}
	}

	session, err := store.Get(r, sessionIdCookie)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userId, err := usuario.Id()
	if err != nil {
		logInternalError(w, database.ErrNotSaved)
		return
	}

	session.Values["userId"] = userId

	if err := session.Save(r, w); err != nil {
		logInternalError(w, err)
		return
	}

	http.Redirect(w, r, "/logged", http.StatusFound)
}

func logged(w http.ResponseWriter, r *http.Request) {
	usuario, err := currentUser(r)
	if err != nil {
		w.Write([]byte("Não Autorizado."))
		return
	}

	t, err := template.ParseFiles(indexLayout, filepath.Join(viewsDir, "logged.html"))
	if err != nil {
		logInternalError(w, err)
		return
	}

	if err := t.Execute(w, usuario.Nome); err != nil {
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

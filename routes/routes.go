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
	// viewsDir    = filepath.Join("views", "html")
	// indexLayout = filepath.Join(viewsDir, "indexLayout.html")
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

type userError struct {
	error
	user string
}

func newUserError(err error, msg string) *userError {
	if err == nil {
		return nil
	}
	return &userError{
		error: err,
		user:  msg,
	}
}

var errUser = errors.New("user error")

type redirectURL string
type redirectFormFunc func(http.ResponseWriter, *http.Request) (redirectURL, error)

func (u redirectURL) Valid() bool {
	return len(u) > 0 && u[0] == '/'
}

func redirectHtmxFormMiddleware(fun redirectFormFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		hxRequest := r.Header.Get("HX-Request") == "true"

		redirectURL, err := fun(w, r)
		if err != nil {
			w.Header().Add("HX-Retarget", "#error-target")
			w.Header().Add("HX-Reswap", "innerHTML")
			w.WriteHeader(http.StatusUnprocessableEntity)
			if u, ok := err.(*userError); ok {
				if hxRequest {
					_ = views.ErrorBox(u.user).Render(context.Background(), w)
				} else {
					w.Write([]byte(u.user))

				}
				return
			} else {
				_ = views.ErrorBox("Erro interno").Render(context.Background(), w)
				return
			}
		}

		if !redirectURL.Valid() {
			panic("not implemented")
		}

		if hxRequest {
			w.Header().Add("HX-Redirect", string(redirectURL))
		} else {
			http.Redirect(w, r, string(redirectURL), http.StatusFound)
		}
	}
}

func logInPOST(w http.ResponseWriter, r *http.Request) (redirectURL, error) {

	email := r.FormValue("email")
	if email == "" {
		return "", formValueMissing("email", w, r)
	}

	senha := r.FormValue("senha")
	if senha == "" {
		return "", formValueMissing("senha", w, r)
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

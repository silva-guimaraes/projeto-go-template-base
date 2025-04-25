package auth

import (
	"fmt"
	"foobar/database"
	"foobar/routes/logging"
	"net/http"
	"os"

	"github.com/gorilla/sessions"
	"github.com/joho/godotenv" // Carrega .env
)

var Store *sessions.CookieStore

const SessionIdCookie = "session_id"

func init() {

	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	key := os.Getenv("SESSION_STORE_KEY")

	if key == "" {
		panic(fmt.Errorf("chave da sessão de cookies não foi definida."))
	}

	Store = sessions.NewCookieStore([]byte(key))
}

func CurrentUser(r *http.Request) (*database.Usuario, error) {
	session, err := Store.Get(r, SessionIdCookie)
	if err != nil {
		logging.Error(err)
		return nil, err
	}

	id, ok := session.Values["userId"].(int)
	if !ok {
		return nil, fmt.Errorf("not ok")
	}

	tx := database.MustBeginTx()
	defer tx.Rollback()

	usuario, err := database.FindUsuarioByID(tx, id)
	if err != nil {
		return nil, err
	}

	return usuario, nil
}

func redirectUnauthorized(w http.ResponseWriter, r *http.Request) {
	session, err := Store.Get(r, SessionIdCookie)
	if err != nil {
		logging.InternalError(w, err)
		return
	}
	session.AddFlash("Autenticação necessária.")
	err = session.Save(r, w)
	if err != nil {
		logging.InternalError(w, err)
		return
	}
	http.Redirect(w, r, "/index.html", http.StatusFound)
}

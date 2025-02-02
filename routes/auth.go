package routes

import (
	"fmt"
	"foobar/database"
	"net/http"
	"os"

	"github.com/gorilla/sessions"
	"github.com/joho/godotenv" // Carrega .env
)

var store *sessions.CookieStore

const sessionIdCookie = "session_id"

func init() {

	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	key := os.Getenv("SESSION_STORE_KEY")

	if key == "" {
		panic(fmt.Errorf("chave da sessão de cookies não foi definida."))
	}

	store = sessions.NewCookieStore([]byte(key))
}

func currentUser(r *http.Request) (*database.Usuario, error) {
	session, err := store.Get(r, sessionIdCookie)
	if err != nil {
		logError(err)
		return nil, err
	}

	id, ok := session.Values["userId"].(int)
	if !ok {
		return nil, fmt.Errorf("não autenticado")
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
	session, err := store.Get(r, sessionIdCookie)
	if err != nil {
		logInternalError(w, err)
		return
	}
	session.AddFlash("Autenticação necessária.")
	err = session.Save(r, w)
	if err != nil {
		logInternalError(w, err)
		return
	}
	http.Redirect(w, r, "/index.html", http.StatusFound)
}

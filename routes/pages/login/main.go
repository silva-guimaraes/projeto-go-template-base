package login

import (
	"database/sql"
	"errors"
	"foobar/database"
	"foobar/model"
	"foobar/routes/auth"
	"foobar/routes/logging"
	"net/http"

	"github.com/a-h/templ"
	"golang.org/x/crypto/bcrypt"
)

func Login(w http.ResponseWriter, r *http.Request, usuario *database.Usuario) (templ.Component, error) {

	userLogged := usuario != nil

	if userLogged {
		return nil, model.NewUserRedirect("/logged")
	} else {
		return login(), nil
	}
}

func LogInPOST(w http.ResponseWriter, r *http.Request) (model.RedirectURL, error) {

	email := r.FormValue("email")
	if email == "" {
		return "", logging.FormValueMissing("email", r)
	}

	senha := r.FormValue("senha")
	if senha == "" {
		return "", logging.FormValueMissing("senha", r)
	}

	tx := database.MustBeginTx()
	defer tx.Rollback()

	usuario, err := database.VerifyUser(tx, email, senha)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) || errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return "", model.NewUserError("senha ou email não conferem")
		} else {
			return "", err
		}
	}

	session, err := auth.Store.Get(r, auth.SessionIdCookie)
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

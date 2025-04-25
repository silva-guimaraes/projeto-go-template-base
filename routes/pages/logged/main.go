package logged

import (
	"foobar/database"
	"foobar/model"
	"net/http"

	"github.com/a-h/templ"
)

func Logged(_ http.ResponseWriter, _ *http.Request, usuario *database.Usuario) (templ.Component, error) {
	if usuario == nil {
		return nil, model.NewUserRedirect("/")
	}
	return logged(usuario), nil
}

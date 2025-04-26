package logged

import (
	"foobar/model"
	"net/http"

	"github.com/a-h/templ"
)

func Logged(_ http.ResponseWriter, _ *http.Request, usuario *model.User) (templ.Component, error) {
	if usuario == nil {
		return nil, model.NewUserRedirect("/")
	}
	return logged(usuario), nil
}

package sobre

import (
	"foobar/database"
	"net/http"

	"github.com/a-h/templ"
)

func Sobre(_ http.ResponseWriter, _ *http.Request, _ *database.Usuario) (templ.Component, error) {
	return sobre(), nil
}

package sobre

import (
	"foobar/model"
	"net/http"

	"github.com/a-h/templ"
)

func Sobre(_ http.ResponseWriter, _ *http.Request, _ *model.User) (templ.Component, error) {
	return sobre(), nil
}

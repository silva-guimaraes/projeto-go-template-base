package routes

import (
	"context"
	"errors"
	"foobar/model"
	. "foobar/model"
	"foobar/routes/auth"
	"foobar/routes/logging"
	"foobar/views"
	"net/http"

	"github.com/a-h/templ"
)

type postRouteRedirectFunc func(http.ResponseWriter, *http.Request) (RedirectURL, error)

type routeGetFunc func(http.ResponseWriter, *http.Request, *model.User) (templ.Component, error)

var (
	errUser = errors.New("user error")
)

func getRouteMiddleware(fun routeGetFunc) http.HandlerFunc {
	// TODO: handle boosted HTMX requests gracefully
	return func(w http.ResponseWriter, r *http.Request) {
		usuario, _ := auth.CurrentUser(r)
		component, err := fun(w, r, usuario)
		if err != nil {
			switch u := err.(type) {
			case *UserRedirect:
				http.Redirect(w, r, u.String(), http.StatusFound)
				return

			default:
				logging.InternalError(w, err)
				return
			}
		} else {
			err = component.Render(context.Background(), w)
			if err != nil {
				logging.InternalError(w, err)
				return
			}
		}
	}
}

func htmxFormRedirectMiddleware(fun postRouteRedirectFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		hxRequest := r.Header.Get("HX-Request") == "true"

		redirectURL, err := fun(w, r)
		if err != nil {
			w.Header().Add("HX-Retarget", "#error-target")
			w.Header().Add("HX-Reswap", "innerHTML")
			w.WriteHeader(http.StatusUnprocessableEntity)
			if u, ok := err.(*UserError); ok {
				if hxRequest {
					_ = views.ErrorBox(u.Error()).Render(context.Background(), w)
				} else {
					w.Write([]byte(u.Error()))

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

package routes

import (
	"context"
	"errors"
	"foobar/views"
	"net/http"
)

type userError struct {
	error
	user string
}

type redirectURL string
type redirectFormFunc func(http.ResponseWriter, *http.Request) (redirectURL, error)

var (
	errUser = errors.New("user error")
)

func newUserError(err error, msg string) *userError {
	if err == nil {
		return nil
	}
	return &userError{
		error: err,
		user:  msg,
	}
}

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

// func indexBoosted(w http.ResponseWriter, r *http.Request, child templ.Component) error {
// 	return renderBoosted(w, r, views.IndexLayout(), child)
// }
//
// func renderBoosted(w http.ResponseWriter, r *http.Request, parent, child templ.Component) error {
// 	isBoosted := r.Header.Get("HX-Boosted") == "true"
//
// 	if isBoosted {
// 		return child.Render(context.Background(), w)
// 	} else {
// 		ctx := templ.WithChildren(context.Background(), child)
// 		return parent.Render(ctx, w)
// 	}
// }

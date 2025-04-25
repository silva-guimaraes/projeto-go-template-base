package model

type RedirectURL string

func (u RedirectURL) Valid() bool {
	return len(u) > 0 && u[0] == '/'
}

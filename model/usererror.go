package model

type UserRedirect struct {
	error
	where string
}

func (u *UserRedirect) String() string {
	return u.where
}

type UserError struct {
	error
	user string
}

func NewUserError(msg string) *UserError {
	return &UserError{
		user: msg,
	}
}

func NewUserRedirect(where string) *UserRedirect {
	return &UserRedirect{
		where: where,
	}
}

func (u *UserError) Error() string {
	return u.user
}

func (u *UserError) String() string {
	return u.user
}

package user

import "errors"

type ErrNoProofKey struct{ err error }

func newErrNoProofKey() ErrNoProofKey {
	return ErrNoProofKey{err: errors.New("expired proof key")}
}

func (e ErrNoProofKey) Error() string {
	return e.err.Error()
}

type ErrAccountExisted struct{ err error }

func newErrAccountExisted() ErrAccountExisted {
	return ErrAccountExisted{err: errors.New("account existed")}
}

func (e ErrAccountExisted) Error() string {
	return e.err.Error()
}

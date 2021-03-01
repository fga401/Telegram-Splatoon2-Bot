package user

import "errors"

// ErrNoProofKey identifies the error that no proof key is found in cache.
type ErrNoProofKey struct{ err error }

func newErrNoProofKey() ErrNoProofKey {
	return ErrNoProofKey{err: errors.New("expired proof key")}
}

func (e ErrNoProofKey) Error() string {
	return e.err.Error()
}

// Is checks if an error is ErrNoProofKey.
func (e *ErrNoProofKey) Is(err error) bool {
	_, ok := err.(*ErrNoProofKey)
	return ok
}

// ErrAccountExisted identifies the error that the account is existed in database.
type ErrAccountExisted struct{ err error }

func newErrAccountExisted() ErrAccountExisted {
	return ErrAccountExisted{err: errors.New("account existed")}
}

func (e ErrAccountExisted) Error() string {
	return e.err.Error()
}

// Is checks if an error is ErrAccountExisted.
func (e *ErrAccountExisted) Is(err error) bool {
	_, ok := err.(*ErrAccountExisted)
	return ok
}

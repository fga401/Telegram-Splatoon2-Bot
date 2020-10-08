package service

import "github.com/pkg/errors"

var (
	// read only
	admins = make(map[int64]struct{})
	// read & write?
	// allowPolling map[int64]struct{}
	// isBlock map[int64]struct{}
)

func loadUsers() {
	adminsList, err := UserTable.LoadAdmin()
	if err != nil {
		panic(errors.Wrap(err, "can't load admins"))
	}
	for _, id := range adminsList {
		admins[id] = struct{}{}
	}
	// todo: load block and polling
}


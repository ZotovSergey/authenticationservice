package myProfilesDB

import (
	"errors"
)

// Сообщения об ошибках при работе с БД
var dbDumpFailErr error = errors.New("failed to dump database")
var noProfileErr error = errors.New("no such profile")
var profileExistsErr error = errors.New("such profile is already exists")
var incorrectPasswordErr error = errors.New("incorrect password")

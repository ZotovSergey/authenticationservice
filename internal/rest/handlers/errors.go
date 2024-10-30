package handlers

import (
	"errors"
)

// Сообщения об ошибках при работе с
var userIsNotAdminErr error = errors.New("access error: authorized user is not an admin")
var canNotEditProfileErr error = errors.New("access error: authorized user can not edit this profile")
var canNotEditPasswordErr error = errors.New("access error: authorized user can not edit this profile password")
var canNotRemoveOwnProfileErr error = errors.New("access error: user tried to remove own profile")
var canNotRemoveOwnProfileFromAdminsErr error = errors.New("access error: user tried to remove own profile from admins list")
var unauthorizedRequestErr error = errors.New("access denied: attempt to authorize unauthorized user")

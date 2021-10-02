package model

import (
	"errors"
)

var (
	// ErrCode is a config or an internal error
	ErrCode = errors.New("Case statement in code is not correct.")
	// ErrNoResult is a not results error
	ErrNoResult = errors.New("Result not found.")
	// ErrDuplicateEntity is a database dupicate record error
	ErrDuplicateEntity = errors.New("Duplicate record found.")
	// ErrUnavailable is a database not available error
	ErrUnavailable = errors.New("Database is unavailable.")
	// ErrUnauthorized is a permissions violation
	ErrUnauthorized = errors.New("User does not have permission to perform this operation.")
	// ErrEntityCreation is database error when fail to create record.
	ErrEntityCreation = errors.New("Could not create record.")
	// ErrEntityDeletion is database error when fail to delete record.
	ErrEntityDeletion = errors.New("Could not delete record.")
	// ErrEntityUpdate is database error when fail to update record.
	ErrEntityUpdate = errors.New("Could not update record.")
)

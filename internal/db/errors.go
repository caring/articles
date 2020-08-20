package db

import "errors"

var (
	// ErrNoRows occurs when no records were found
	ErrNoRows = errors.New("no rows")
	// ErrNoRowsAffected occurs when no rows were updated
	ErrNoRowsAffected = errors.New("no rows affected")
	// ErrNotFound when a specific reqcord was not found
	ErrNotFound = errors.New("the record you are attempting to update is not found")
)

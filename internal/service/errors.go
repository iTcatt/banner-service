package service

import "errors"

var (
	ErrNoPermission  = errors.New("permission denied")
	ErrAlreadyExists = errors.New("banner with tag_id and feature_id already exists")
)

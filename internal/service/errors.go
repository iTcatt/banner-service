package service

import "errors"

var (
	ErrAlreadyExists = errors.New("banner with tag_id and feature_id already exists")
)

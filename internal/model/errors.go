package model

import "errors"

var (
	ErrTaskNotFound  = errors.New("task not found")
	ErrTitleRequired = errors.New("title is required")
	ErrInvalidStatus = errors.New("invalid status")
)

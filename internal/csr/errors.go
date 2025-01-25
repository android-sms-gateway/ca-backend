package csr

import "errors"

var (
	ErrInvalidCSR       = errors.New("invalid csr")
	ErrCSRNotFound      = errors.New("csr not found")
	ErrCSRAlreadyExists = errors.New("csr already exists")
)

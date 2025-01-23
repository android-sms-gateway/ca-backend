package csr

import "errors"

var (
	ErrCSRNotFound      = errors.New("csr not found")
	ErrCSRAlreadyExists = errors.New("csr already exists")
)

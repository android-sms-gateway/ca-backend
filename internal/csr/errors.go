package csr

import "errors"

var (
	ErrCACertInvalid    = errors.New("invalid ca certificate")
	ErrCSRInvalid       = errors.New("invalid csr")
	ErrCSRNotFound      = errors.New("csr not found")
	ErrCSRAlreadyExists = errors.New("csr already exists")
)

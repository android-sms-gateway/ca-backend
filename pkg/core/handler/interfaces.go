package handler

type Validatable interface {
	Validate() error
}

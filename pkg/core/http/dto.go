package http

type Error struct {
	Message string `json:"message"`
	Details any    `json:"details,omitempty"`
}

type ErrorResponse struct {
	Error Error `json:"error"`
}

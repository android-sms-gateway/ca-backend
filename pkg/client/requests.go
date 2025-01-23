package client

// PostCSRRequest represents a request to post a Certificate Signing Request (CSR).
type PostCSRRequest struct {
	// Content contains the CSR content and is required.
	Content string `json:"content" validate:"required"`
	// Metadata includes additional metadata related to the CSR.
	Metadata map[string]string `json:"metadata"`
}

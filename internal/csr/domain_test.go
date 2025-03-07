package csr_test

import (
	"testing"

	"github.com/android-sms-gateway/ca/internal/csr"
	"github.com/android-sms-gateway/client-go/ca"
	"github.com/go-playground/assert/v2"
)

func TestNewCSR(t *testing.T) {
	tests := []struct {
		name     string
		csrType  ca.CSRType
		content  string
		metadata map[string]string
	}{
		{
			name:     "With type provided",
			csrType:  ca.CSRTypePrivateServer,
			content:  "test-content",
			metadata: map[string]string{"key": "value"},
		},
		{
			name:     "With empty type",
			csrType:  "",
			content:  "test-content",
			metadata: map[string]string{"key": "value"},
		},
		{
			name:     "With empty metadata",
			csrType:  ca.CSRTypePrivateServer,
			content:  "test-content",
			metadata: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := csr.NewCSR(tt.csrType, tt.content, tt.metadata)

			// Check type (should default to webhook if empty)
			expectedType := tt.csrType
			if expectedType == "" {
				expectedType = ca.CSRTypeWebhook
			}
			assert.Equal(t, expectedType, got.Type())

			// Check content and metadata
			assert.Equal(t, tt.content, got.Content())
			assert.Equal(t, tt.metadata, got.Metadata())
		})
	}
}

func TestCSR_Getters(t *testing.T) {
	csrType := ca.CSRType("private_server")
	content := "test-content"
	metadata := map[string]string{"key": "value"}

	testCSR := csr.NewCSR(csrType, content, metadata)

	assert.Equal(t, csrType, testCSR.Type())
	assert.Equal(t, content, testCSR.Content())
	assert.Equal(t, metadata, testCSR.Metadata())
}

func TestNewCSRStatus(t *testing.T) {
	id := "test-id"
	csrType := ca.CSRType("private_server")
	content := "test-content"
	metadata := map[string]string{"key": "value"}
	status := ca.CSRStatus("pending")
	certificate := "test-certificate"
	reason := "test-reason"

	csrStatus := csr.NewCSRStatus(id, csrType, content, metadata, status, certificate, reason)

	// Verify through public getters
	assert.Equal(t, id, csrStatus.ID())
	assert.Equal(t, csrType, csrStatus.Type())
	assert.Equal(t, content, csrStatus.Content())
	assert.Equal(t, metadata, csrStatus.Metadata())
	assert.Equal(t, status, csrStatus.Status())
	assert.Equal(t, certificate, csrStatus.Certificate())
}

func TestCSRStatus_Getters(t *testing.T) {
	id := "test-id"
	csrType := ca.CSRType("private_server")
	content := "test-content"
	metadata := map[string]string{"key": "value"}
	status := ca.CSRStatus("approved")
	certificate := "test-certificate"
	reason := "test-reason"

	csrStatus := csr.NewCSRStatus(id, csrType, content, metadata, status, certificate, reason)

	assert.Equal(t, id, csrStatus.ID())
	assert.Equal(t, status, csrStatus.Status())
	assert.Equal(t, certificate, csrStatus.Certificate())
}

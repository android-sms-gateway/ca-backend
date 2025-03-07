package api

import (
	"github.com/android-sms-gateway/ca/internal/csr"
	"github.com/android-sms-gateway/client-go/ca"
)

func csrStatusToResponse(status csr.CSRStatus) ca.PostCSRResponse {
	return ca.PostCSRResponse{
		RequestID:   status.ID(),
		Type:        status.Type(),
		Status:      status.Status(),
		Message:     status.Status().Description(),
		Certificate: status.Certificate(),
	}
}

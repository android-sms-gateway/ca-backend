package csr

import (
	"encoding/json"

	"github.com/android-sms-gateway/ca/pkg/client"
)

type CSR struct {
	content  string
	metadata map[string]string
}

func NewCSR(content string, metadata map[string]string) CSR {
	return CSR{
		content:  content,
		metadata: metadata,
	}
}

func (c CSR) Content() string {
	return c.content
}

func (c CSR) Metadata() map[string]string {
	return c.metadata
}

func (c CSR) toMap() map[string]string {
	metadata := "{}"
	if len(c.metadata) > 0 {
		b, _ := json.Marshal(c.metadata)
		metadata = string(b)
	}

	return map[string]string{
		"content":  c.content,
		"metadata": metadata,
	}
}

type CSRStatus struct {
	id          string
	status      client.CSRStatus
	certificate string
	reason      string
}

func NewCSRStatus(id string, status client.CSRStatus, certificate string, reason string) CSRStatus {
	return CSRStatus{
		id:          id,
		status:      status,
		certificate: certificate,
		reason:      reason,
	}
}

func (c CSRStatus) ID() string {
	return c.id
}

func (c CSRStatus) Status() client.CSRStatus {
	return c.status
}

func (c CSRStatus) Certificate() string {
	return c.certificate
}

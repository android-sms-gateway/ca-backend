package csr

import (
	"encoding/json"

	"github.com/android-sms-gateway/client-go/ca"
)

type CSR struct {
	csrType  ca.CSRType
	content  string
	metadata map[string]string
}

func NewCSR(csrType ca.CSRType, content string, metadata map[string]string) CSR {
	if csrType == "" {
		csrType = ca.CSRTypeWebhook
	}

	return CSR{
		csrType:  csrType,
		content:  content,
		metadata: metadata,
	}
}

func (c CSR) Type() ca.CSRType {
	return c.csrType
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
		"type":     string(c.csrType),
		"content":  c.content,
		"metadata": metadata,
	}
}

type Status struct {
	CSR

	id          string
	status      ca.CSRStatus
	certificate string
	reason      string
}

func NewCSRStatus(
	id string,
	csrType ca.CSRType,
	content string,
	metadata map[string]string,
	status ca.CSRStatus,
	certificate string,
	reason string,
) Status {
	return Status{
		CSR: CSR{
			csrType:  csrType,
			content:  content,
			metadata: metadata,
		},
		id:          id,
		status:      status,
		certificate: certificate,
		reason:      reason,
	}
}

func (c Status) ID() string {
	return c.id
}

func (c Status) Status() ca.CSRStatus {
	return c.status
}

func (c Status) Certificate() string {
	return c.certificate
}

type csrID string

func (c csrID) Bytes() []byte {
	return []byte(c)
}

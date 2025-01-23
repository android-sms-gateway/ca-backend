package csr

import "github.com/android-sms-gateway/ca/pkg/client"

type CSR struct {
	Content  string
	Metadata map[string]string
}

type CSRStatus struct {
	Status      client.CSRStatus
	Certificate string
}

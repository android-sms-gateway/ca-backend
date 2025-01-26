package csr

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
)

func LoadCA(certPEM, keyPEM []byte) (*x509.Certificate, any, error) {
	// Load CA certificate
	block, _ := pem.Decode(certPEM)
	if block == nil || block.Type != "CERTIFICATE" {
		return nil, nil, fmt.Errorf("invalid CA certificate")
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse CA certificate: %w", err)
	}

	// Load CA private key
	block, _ = pem.Decode(keyPEM)
	if block == nil || (block.Type != "RSA PRIVATE KEY" && block.Type != "EC PRIVATE KEY" && block.Type != "PRIVATE KEY") {
		return nil, nil, fmt.Errorf("invalid CA private key")
	}
	priv, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse CA private key: %w", err)
	}

	return cert, priv, nil
}

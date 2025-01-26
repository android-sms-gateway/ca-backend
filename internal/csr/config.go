package csr

import "time"

type Config struct {
	CACert []byte
	CAKey  []byte
	TTL    time.Duration
}

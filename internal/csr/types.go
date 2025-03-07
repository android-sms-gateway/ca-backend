package csr

import "github.com/android-sms-gateway/client-go/ca"

type SerialNumberPrefix uint8

const (
	PrefixWebhooks       SerialNumberPrefix = 1
	PrefixPrivateServers SerialNumberPrefix = 2
)

var csrTypeToPrefix = map[ca.CSRType]SerialNumberPrefix{
	ca.CSRTypeWebhook:       PrefixWebhooks,
	ca.CSRTypePrivateServer: PrefixPrivateServers,
}

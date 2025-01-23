package main

import "github.com/android-sms-gateway/ca/internal"

//go:generate swag init --parseDependency -g ./main.go -o ./api

//	@title			SMS Gate Certificate Authority API
//	@version		{{VERSION}}
//	@description	Provides methods to manage certificates

//	@contact.name	SMSGate Support
//	@contact.email	support@sms-gate.app

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@host		ca.sms-gate.app
//	@BasePath	/api/v1
//  @schemes	https

func main() {
	internal.Run()
}

package version

import "strconv"

const notSet string = "not set"

// These variables are populated at build time using ldflags.
// Example: go build -ldflags "-X github.com/android-sms-gateway/ca/internal/version.AppVersion=0.1 -X github.com/android-sms-gateway/ca/internal/version.AppRelease=123"
var (
	AppVersion = notSet
	AppRelease = notSet
)

func AppReleaseID() int {
	id, _ := strconv.Atoi(AppRelease)

	return id
}

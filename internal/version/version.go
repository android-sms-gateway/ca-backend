package version

import "strconv"

const notSet string = "not set"

//nolint:gochecknoglobals // These variables are populated at build time using ldflags.
var (
	AppVersion = notSet
	AppRelease = notSet
)

func AppReleaseID() int {
	id, _ := strconv.Atoi(AppRelease)

	return id
}

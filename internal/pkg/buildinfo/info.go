package buildinfo

// These variables are set during build time via -ldflags
var (
	Version   string
	BuildTime string
)

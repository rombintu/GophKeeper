package common

import "fmt"

func ifEmptyOpt(opt string) string {
	if opt == "" {
		return "N/A"
	}
	return opt
}

func Version(buildDate, buildVersion, buildCommit string) string {
	return fmt.Sprintf("Build date: %s\nBuild version: %s\nBuild commit: %s",
		ifEmptyOpt(buildDate),
		ifEmptyOpt(buildVersion),
		ifEmptyOpt(buildCommit),
	)
}

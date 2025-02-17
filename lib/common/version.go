package common

import "log/slog"

func ifEmptyOpt(opt string) string {
	if opt == "" {
		return "N/A"
	}
	return opt
}

func Version(buildDate, buildVersion, buildCommit, binary string) {
	slog.Info(
		"Init", slog.String("Binary", ifEmptyOpt(binary)),
		slog.String("Build version", ifEmptyOpt(buildVersion)),
		slog.String("Build date", ifEmptyOpt(buildDate)),
		slog.String("Build commit", ifEmptyOpt(buildCommit)),
	)
}

package cli

type Profile struct {
	Email string
	Token string
}

// master, err := crypto.LoadPrivateKey(*privateKeyPath)
// if err != nil {
// 	slog.Error("load master key", slog.String("error", err.Error()))
// 	os.Exit(0)
// }

// user, err := crypto.GetProfile(master)
// if err != nil {
// 	slog.Error("load get profile", slog.String("error", err.Error()))
// 	os.Exit(0)
// }

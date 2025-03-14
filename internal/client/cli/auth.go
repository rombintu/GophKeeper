package cli

// func auth() {
// 	// Создание пула соединений, из него создаются клиенты
// 	connPool := connections.NewConnPool()
// 	defer connPool.CleanUp()
// 	clientPool := client.NewClientPool(connPool)
// 	authClient, err := clientPool.NewAuthClient(*addressAuth)
// 	if err != nil {
// 		slog.Error("get connection to auth service", slog.String("error", err.Error()))
// 		os.Exit(0)
// 	}
// }

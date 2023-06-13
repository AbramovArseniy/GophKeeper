package main

import (
	"context"
	"log"

	// "os"

	"github.com/AbramovArseniy/GophKeeper/internal/client"
	"github.com/AbramovArseniy/GophKeeper/internal/client/utils/config"

	// "github.com/AbramovArseniy/GophKeeper/internal/server/handlers"
	// "github.com/AbramovArseniy/GophKeeper/internal/server/utils/config"
	"google.golang.org/grpc/metadata"
)

// func main() {
// 	context := context.Background()
// 	enableTLS := os.Getenv("ENABLE_TLS") == "true"
// 	address := os.Getenv("SERVER_ADDRESS")
// 	address = "127.0.0.1:9000"
// 	enableTLS = false
// 	cfg := config.SetServerParams()
// 	auth := handlers.NewAuth(context, database, cfg.JWTSecret)

// 	client := client.NewCLI()
// 	err := client.StartCLI(context)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// }

func main() {
	var md metadata.MD
	metadate := metadata.New(map[string]string{})
	ctx := metadata.NewOutgoingContext(context.Background(), metadate)
	cfg := config.SetClientParams()
	action, err := client.NewAction(cfg.ServerAddr, &md)
	if err != nil {
		log.Fatal("Failed connect to server")
	}

	client := client.NewCLI(action)
	err = client.StartCLI(ctx)
	if err != nil {
		log.Fatal(err)
	}
}

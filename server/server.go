package server

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/stodis/stodis/api/protobuf/services/fileservice"
	"github.com/stodis/stodis/internal/discord"
	"github.com/stodis/stodis/internal/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func Serve() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "50051"
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Create a new gRPC server
	s := grpc.NewServer()
	token := os.Getenv("DISCORD_BOT_TOKEN")
	disc, err := discord.NewBot(discord.NewBotConfig(token))
	if err != nil {
		log.Fatalf("failed to create discord bot: %v", err)
	}

	_ = discord.NewDiscordFileService(disc, "1278013883973632071")

	// Register the server with the gRPC server
	fileservice.RegisterUploadFileServer(s, &service.Server{})
	reflection.Register(s)

	// Start the server
	fmt.Println("Server is running on port", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

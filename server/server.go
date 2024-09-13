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

var DiscordToken []string = []string{
	
}

func createDiscordBots() ([]*discord.Bot, error) {
	var bots []*discord.Bot
	for _, token := range DiscordToken {
		disc, err := discord.NewBot(discord.NewBotConfig(token))
		if err != nil {
			log.Fatalf("failed to create discord bot: %v", err)
			return nil, err
		}
		bots = append(bots, disc)
	}
	return bots, nil
}

func createDiscordHandler() ([]service.FileService, error) {
	discordBots, err := createDiscordBots()
	if err != nil {
		log.Fatalf("failed to create discord bots: %v", err)
		return nil, err
	}
	var discHandlers []service.FileService
	for _, bot := range discordBots {
		discHandlers = append(discHandlers, discord.NewDiscordFileService(bot, "1278013883973632071"))
	}
	return discHandlers, nil
}

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
	// token := os.Getenv("DISCORD_BOT_TOKEN")
	// disc, err := discord.NewBot(discord.NewBotConfig(token))
	if err != nil {
		log.Fatalf("failed to create discord bot: %v", err)
	}
	// discordBots := createDiscordBots()
	// discodHandler := discord.NewDiscordFileService(disc, "1278013883973632071")
	discordHandler, err := createDiscordHandler()
	if err != nil {
		log.Fatalf("failed to create discord handler: %v", err)
	}
	// Register the server with the gRPC server
	fileservice.RegisterUploadFileServer(s, service.NewServer(discordHandler))
	reflection.Register(s)

	// Start the server
	fmt.Println("Server is running on port", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

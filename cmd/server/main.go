package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"os"

	"github.com/stodis/stodis/api/protobuf/services/fileservice"
	"github.com/stodis/stodis/internal/discord"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct {
	fileservice.UnimplementedUploadFileServer
}

// Implement the CreateFile RPC method
func (s *server) CreateFile(ctx context.Context, req *fileservice.CreateFileRequest) (*fileservice.CreateFileResponse, error) {
	fmt.Printf("Received CreateFile request: %v\n", req)
	// Here you would normally handle the request and perform your logic
	return &fileservice.CreateFileResponse{Uuid: "generated-uuid"}, nil
}

// Implement the UploadFile RPC method
func (s *server) UploadFile(stream fileservice.UploadFile_UploadFileServer) error {
	var fileChunks []*fileservice.FileData
	for {
		chunk, err := stream.Recv()
		if err == io.EOF {
			// End of file stream
			break
		}
		if err != nil {
			return err
		}
		fileChunks = append(fileChunks, chunk)
	}

	fmt.Printf("Received file with %d chunks\n", len(fileChunks))
	// Here you would normally handle the file chunks and perform your logic
	return stream.SendAndClose(&fileservice.FileUploadResponse{Message: "File uploaded successfully", Success: true})
}

func main() {
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
	fileservice.RegisterUploadFileServer(s, &server{})
	reflection.Register(s)

	// Start the server
	fmt.Println("Server is running on port", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

package main

import (
	"context"
	"stodis/api/protobuf/services/fileservice" // Import the package where your generated files are located
	"fmt"
	"io"
	"log"
	"net"

	"google.golang.org/grpc"
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
	var fileChunks []fileservice.FileData
	for {
		chunk, err := stream.Recv()
		if err == io.EOF {
			// End of file stream
			break
		}
		if err != nil {
			return err
		}
		fileChunks = append(fileChunks, *chunk)
	}

	fmt.Printf("Received file with %d chunks\n", len(fileChunks))
	// Here you would normally handle the file chunks and perform your logic
	return stream.SendAndClose(&fileservice.FileUploadResponse{Message: "File uploaded successfully", Success: true})
}

func main() {
	// Set up a connection to the server.
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Create a new gRPC server
	s := grpc.NewServer()

	// Register the server with the gRPC server
	fileservice.RegisterUploadFileServer(s, &server{})

	// Start the server
	fmt.Println("Server is running on port :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

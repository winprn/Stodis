package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/stodis/stodis/api/protobuf/services/fileservice"
	"google.golang.org/grpc"
)

func main() {
	// Connect to the gRPC server
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to connect to server: %v", err)
	}
	defer conn.Close()

	client := fileservice.NewUploadFileClient(conn)

	// Create the file on the server by calling CreateFile
	createFileReq := &fileservice.CreateFileRequest{
		FileName: "100MB_file.bin",
		FileSize: 1024,                       // Size of the file in bytes
		FileType: fileservice.FileType_image, // Specify the file type (image, document, media)
	}
	createFileResp, err := client.CreateFile(context.Background(), createFileReq)
	if err != nil {
		log.Fatalf("failed to create file: %v", err)
	}
	fmt.Printf("File created with UUID: %s\n", createFileResp.Uuid)

	// Open the file to be uploaded
	filePath := "100MB_file.bin"
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("failed to open file: %v", err)
	}
	defer file.Close()

	// Prepare file for upload in chunks
	stream, err := client.UploadFile(context.Background())
	if err != nil {
		log.Fatalf("failed to start file upload: %v", err)
	}

	buffer := make([]byte, 64*1024) // 64 KB chunk size
	chunkCounter := int32(0)        // Keep track of the chunk number

	for {
		n, err := file.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("failed to read file: %v", err)
		}

		// Send each file chunk
		chunkCounter++
		req := &fileservice.FileData{
			FileId:  createFileResp.Uuid,
			Chunk:   buffer[:n],
			ChunkTh: chunkCounter, // Add chunk index number
		}
		if err := stream.Send(req); err != nil {
			log.Fatalf("failed to send chunk: %v", err)
		}
	}

	// Close the stream and receive the response
	resp, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("failed to receive response: %v", err)
	}

	fmt.Printf("File upload response: %s, success: %v\n", resp.Message, resp.Success)
}

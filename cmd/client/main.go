package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/stodis/stodis/api/protobuf/services/fileservice"
	"google.golang.org/grpc"
)

func uploadFile(fileName string, client *fileservice.UploadFileClient) error {
	createFileReq := &fileservice.CreateFileRequest{
		FileName: fileName,
		FileSize: 1024,                       // Size of the file in bytes
		FileType: fileservice.FileType_image, // Specify the file type (image, document, media)
	}
	createFileResp, err := (*client).CreateFile(context.Background(), createFileReq)
	if err != nil {
		log.Fatalf("failed to create file: %v", err)
		return err
	}
	fmt.Printf("File created with UUID: %s\n", createFileResp.Uuid)
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("failed to open file: %v", err)
		return err
	}
	defer file.Close()
	stream, err := (*client).UploadFile(context.Background())
	if err != nil {
		log.Fatalf("failed to start file upload: %v", err)
		return err
	}
	buffer := make([]byte, 64*1024) // 64 KB chunk size
	chunkCounter := int32(0)        // Keep track of the chunk number
	// return nil
	startTime := time.Now()
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
	endTime := time.Now()
	fmt.Printf("Time taken to upload file: %v second\n", endTime.Sub(startTime))
	return nil
}

func main() {
	// Connect to the gRPC server
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to connect to server: %v", err)
	}
	defer conn.Close()

	client := fileservice.NewUploadFileClient(conn)
	// uploadFile("test.txt", &client)
	uploadFile("100MB_file_2.bin", &client)
}

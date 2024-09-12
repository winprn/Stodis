package service

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"time"

	"github.com/google/uuid"
	"github.com/stodis/stodis/api/protobuf/services/fileservice"
)

const (
	chunkSize = 25 * 1024 * 1024 // 25 MB
)

type Server struct {
	fileservice.UnimplementedUploadFileServer

	discordService FileService
}

func NewServer(discordService FileService) *Server {
	return &Server{discordService: discordService}
}

// Implement the CreateFile RPC method
func (s *Server) CreateFile(ctx context.Context, req *fileservice.CreateFileRequest) (*fileservice.CreateFileResponse, error) {
	fmt.Printf("Received CreateFile request: %v\n", req)
	// Here you would normally handle the request and perform your logic
	uuid := uuid.New().String()
	return &fileservice.CreateFileResponse{Uuid: uuid}, nil
}

func storeChunk(chunks *[]bytes.Buffer, buffer *bytes.Buffer, chunk []byte) {
	startIndex := 0
	for {
		writeSize := min(len(chunk), chunkSize-buffer.Len())
		buffer.Write(chunk[startIndex : startIndex+writeSize])
		startIndex += writeSize
		if buffer.Len() == chunkSize {
			*chunks = append(*chunks, *buffer)
			buffer.Reset()
		}
		if startIndex == len(chunk) {
			break
		}
	}
}

// Implement the UploadFile RPC method
func (s *Server) UploadFile(stream fileservice.UploadFile_UploadFileServer) error {
	var fileChunks []*fileservice.FileData
	var chunks []bytes.Buffer
	var buffer bytes.Buffer
	for {
		chunk, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		fileChunks = append(fileChunks, chunk)
		storeChunk(&chunks, &buffer, chunk.GetChunk())
	}
	if buffer.Len() > 0 {
		chunks = append(chunks, buffer)
	}
	startTime := time.Now()
	for index, chunk := range chunks {
		fileName := fmt.Sprintf("%s-%d", fileChunks[0].GetFileId(), index)
		if _, err := s.discordService.UploadFile(chunk.Bytes(), fileName); err != nil {
			return err
		}
	}
	endTime := time.Now()
	fmt.Printf("Time taken to upload file: %v second\n", endTime.Sub(startTime))
	fmt.Printf("Received file with %d chunks\n", len(fileChunks))
	return stream.SendAndClose(&fileservice.FileUploadResponse{Message: "File uploaded successfully", Success: true})
}

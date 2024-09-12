package service

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/google/uuid"
	"github.com/stodis/stodis/api/protobuf/services/fileservice"
)

const (
	chunkSize = 25 * 1024 * 1024 // 25 MB
)

type Server struct {
	fileservice.UnimplementedUploadFileServer

	fileService FileService
}

func NewServer(fileService FileService) *Server {
	return &Server{fileService: fileService}
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
			if buffer.Len() > 0 {
				*chunks = append(*chunks, *buffer)
				buffer.Reset()
			}
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
			// End of file stream
			break
		}
		if err != nil {
			return err
		}
		fileChunks = append(fileChunks, chunk)
		storeChunk(&chunks, &buffer, chunk.GetChunk())
	}
	fmt.Printf("Number of chunks: %d\n", len(chunks))
	for _, chunk := range chunks {
		if _, err := s.fileService.UploadFile(chunk.Bytes(), "testfile.txt"); err != nil {
			return err
		}
		fmt.Printf("Chunk: %s\n", chunk.String())
	}
	fmt.Printf("Received file with %d chunks\n", len(fileChunks))
	return stream.SendAndClose(&fileservice.FileUploadResponse{Message: "File uploaded successfully", Success: true})
}

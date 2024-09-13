package service

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/stodis/stodis/api/protobuf/services/fileservice"
)

const (
	chunkSize = 25 * 1024 * 1024 // 25 MB
	BotWorker = 2
)

type fileData struct {
	fileId  string
	data    bytes.Buffer
	chunkTh int32
}

type Server struct {
	fileservice.UnimplementedUploadFileServer

	discordService []FileService
	chunks         chan fileData
}

func NewServer(discordService []FileService) *Server {
	server := &Server{
		discordService: discordService,
		chunks:         make(chan fileData, 100),
	}
	for botId := 0; botId < BotWorker; botId++ {
		go server.flush(botId)
	}
	return server
}

// Implement the CreateFile RPC method
func (s *Server) CreateFile(ctx context.Context, req *fileservice.CreateFileRequest) (*fileservice.CreateFileResponse, error) {
	fmt.Printf("Received CreateFile request: %v\n", req)
	uuid := uuid.New().String()
	return &fileservice.CreateFileResponse{Uuid: uuid}, nil
}

func (s *Server) storeChunk(buffer *bytes.Buffer, chunk []byte, fileId string, chunkTh *int32) {
	startIndex := 0
	for {
		writeSize := min(len(chunk), chunkSize-buffer.Len())
		buffer.Write(chunk[startIndex : startIndex+writeSize])
		startIndex += writeSize
		if buffer.Len() == chunkSize {
			data := new(bytes.Buffer)
			io.Copy(data, buffer)
			s.chunks <- fileData{
				data:    *data,
				chunkTh: *chunkTh,
				fileId:  fileId,
			}
			*chunkTh += 1
			buffer.Reset()
		}
		if startIndex == len(chunk) {
			break
		}
	}
}

// Implement the UploadFile RPC method with concurrency
func (s *Server) UploadFile(stream fileservice.UploadFile_UploadFileServer) error {
	var buffer bytes.Buffer
	chunkTh := int32(0)
	startTime := time.Now()
	cnt := 0
	id := "-1"
	for {
		chunk, err := stream.Recv()
		cnt += 1
		if err == io.EOF {
			if buffer.Len() > 0 {
				s.chunks <- fileData{
					fileId:  id,
					data:    buffer,
					chunkTh: chunkTh,
				}
			}
			log.Printf("File upload completed\n")
			break
		}
		if id == "-1" {
			id = chunk.GetFileId()
		} else if id != chunk.GetFileId() {
			return fmt.Errorf("file id mismatch")
		}
		if err != nil {
			return err
		}
		s.storeChunk(&buffer, chunk.GetChunk(), id, &chunkTh)
	}
	endTime := time.Now()
	fmt.Printf("Time taken to upload file: %v\n", endTime.Sub(startTime))

	return stream.SendAndClose(&fileservice.FileUploadResponse{Message: "File uploaded successfully", Success: true})
}

func (s *Server) flush(botId int) (err error) {
	for chunk := range s.chunks {
		fileName := fmt.Sprintf("%s-%d", chunk.fileId, chunk.chunkTh)
		data := chunk.data.Bytes()
		if _, err := s.discordService[botId].UploadFile(data, fileName); err != nil {
			log.Printf("failed to upload file: %v\n", err)
			s.chunks <- chunk
			time.Sleep(1 * time.Second)
		}
	}
	return nil
}

package service

type FileService interface {
	UploadFile(file []byte, name string) (string, error)
}

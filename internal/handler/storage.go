package handler

type StorageHandler interface {
	UploadFile(file []byte, name string) (string, error)
}

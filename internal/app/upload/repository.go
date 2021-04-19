package upload

import "mime/multipart"

type UploadRepository interface {
	InsertPhotos(files []*multipart.FileHeader, photoPath string) ([]string, error)
}

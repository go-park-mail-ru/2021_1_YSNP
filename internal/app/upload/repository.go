package upload

import "mime/multipart"

//go:generate mockgen -destination=./mocks/mock_upload_repo.go -package=mock github.com/go-park-mail-ru/2021_1_YSNP/internal/app/upload UploadRepository

type UploadRepository interface {
	InsertPhotos(files []*multipart.FileHeader, photoPath string) ([]string, error)

	RemovePhoto(imgUrl string) error
	RemovePhotos(imgUrls []string) error
}

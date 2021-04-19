package FileSystem

import (
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/upload"
	"github.com/google/uuid"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

type UploadRepository struct {
}

func NewUploadRepository() upload.UploadRepository {
	return &UploadRepository{}
}

func (ur *UploadRepository) InsertPhotos(files []*multipart.FileHeader, photoPath string) ([]string, error) {
	imgUrls := make(map[string][]string)

	for i := range files {
		file, err := files[i].Open()
		if err != nil {
			return nil, err
		}
		defer file.Close()

		str, err := os.Getwd()
		if err != nil {
			return nil, err
		}

		os.Chdir(photoPath)

		photoID, err := uuid.NewRandom()
		if err != nil {
			return nil, err
		}

		extension := filepath.Ext(files[i].Filename)

		newFile, err := os.OpenFile(photoID.String()+extension, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			return nil, err
		}
		defer newFile.Close()

		os.Chdir(str)

		_, err = io.Copy(newFile, file)
		if err != nil {
			_ = os.Remove(photoID.String() + extension)
			return nil, err
		}

		imgUrls["img"] = append(imgUrls["img"], "/"+photoPath+photoID.String()+extension)
	}

	return imgUrls["img"], nil
}
